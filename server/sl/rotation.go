// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Rotation struct {
	PluginVersion string
	RotationID    types.ID
	IsArchived    bool

	FillerType types.ID
	TaskType   types.ID

	MattermostUserIDs *types.IDSet `json:",omitempty"`
	TaskIDs           *types.IDSet `json:",omitempty"`

	TaskSettings      TaskSettings      `json:",omitempty"`
	FillSettings      FillSettings      `json:",omitempty"`
	AutopilotSettings AutopilotSettings `json:",omitempty"`

	loaded bool
	Users  *Users `json:"-"`
	Tasks  *Tasks `json:"-"`
}

type TaskSettings struct {
	Seq         int           `json:",omitempty"`
	Require     *Needs        `json:",omitempty"`
	Limit       *Needs        `json:",omitempty"`
	Duration    time.Duration `json:",omitempty"`
	Grace       time.Duration `json:",omitempty"`
	Description string        `json:",omitempty"`
}

type FillSettings struct {
	// Beginning marks the time when scheduling starts
	Beginning types.Time

	// Period specifies the recurrence of shifts, but is also used to calculate
	// the user weights when filling tasks.
	Period types.Period

	// Seed is a stable seed to use for rand.
	Seed int64 `json:",omitempty"`

	// Fuzz is the number of periods that gets added to the default doubling
	// duration when calculating user weights.
	Fuzz int64 `json:",omitempty"`
}

type AutopilotSettings struct {
	Create            bool          `json:",omitempty"`
	CreatePrior       time.Duration `json:",omitempty"`
	Schedule          bool          `json:",omitempty"`
	SchedulePrior     time.Duration `json:",omitempty"`
	StartFinish       bool          `json:",omitempty"`
	RemindStart       bool          `json:",omitempty"`
	RemindStartPrior  time.Duration `json:",omitempty"`
	RemindFinish      bool          `json:",omitempty"`
	RemindFinishPrior time.Duration `json:",omitempty"`
}

const (
	TaskTypeTicket = types.ID("ticket")
	TaskTypeShift  = types.ID("shift")
)

func NewRotation() *Rotation {
	r := &Rotation{}
	r.Init()
	return r
}

func (r *Rotation) Init() {
	if r.MattermostUserIDs == nil {
		r.MattermostUserIDs = types.NewIDSet()
	}
	if r.TaskIDs == nil {
		r.TaskIDs = types.NewIDSet()
	}
	if r.TaskSettings.Require == nil {
		r.TaskSettings.Require = NewNeeds()
	}
	if r.TaskSettings.Limit == nil {
		r.TaskSettings.Limit = NewNeeds()
	}
}

func (rotation *Rotation) WithMattermostUserIDs(pool *Users) *Rotation {
	newRotation := *rotation
	newRotation.MattermostUserIDs = types.NewIDSet()
	for _, id := range pool.IDs() {
		newRotation.MattermostUserIDs.Set(id)
	}
	if pool.IsEmpty() {
		pool = NewUsers()
	}
	newRotation.Users = pool
	return &newRotation
}

func (r *Rotation) String() string {
	return r.Name()
}

func (r *Rotation) Name() string {
	return kvstore.NameFromID(r.RotationID)
}

func (r *Rotation) Markdown() md.MD {
	return md.MD(r.Name())
}

func (r *Rotation) MarkdownBullets() md.MD {
	out := md.Markdownf("- **%s**\n", r.Name())
	out += md.Markdownf("  - ID: `%s`.\n", r.RotationID)
	if r.Users != nil {
		out += md.Markdownf("  - Users (%v): %s.\n", r.MattermostUserIDs.Len(), r.Users.MarkdownWithSkills())
	} else {
		out += md.Markdownf("  - Users (%v): %s.\n", r.MattermostUserIDs.Len(), r.MattermostUserIDs.IDs())
	}

	out += md.Markdownf("  - Task settings:\n")
	out += md.Markdownf("    - Task type: **%s**\n", r.TaskType)
	out += md.Markdownf("    - Require: %s\n", r.TaskSettings.Require.Markdown())
	out += md.Markdownf("    - Limit: %v\n", r.TaskSettings.Limit.Markdown())
	out += md.Markdownf("    - Grace: **%v**\n", r.TaskSettings.Grace)

	out += md.Markdownf("  - Fill settings:\n")
	out += md.Markdownf("    - Filler type: **%s**\n", r.FillerType)
	out += md.Markdownf("    - Beginning: **%s**\n", r.FillSettings.Beginning)
	out += md.Markdownf("    - Shift period: **%s**\n", r.FillSettings.Period)
	out += md.Markdownf("    - Fuzz: **%v**\n", r.FillSettings.Fuzz)

	if r.AutopilotSettings.isOn() {
		out += md.Markdownf("  - Autopilot: **on**\n")

		if r.AutopilotSettings.Create {
			out += md.Markdownf("    - Create shifts **%v** prior to start\n", r.AutopilotSettings.CreatePrior)
		}
		if r.AutopilotSettings.Schedule {
			out += md.Markdownf("    - Fill and schedule tasks **%v** prior to start\n", r.AutopilotSettings.SchedulePrior)
		}
		if r.AutopilotSettings.StartFinish {
			out += md.Markdownf("    - Start and finish scheduled tasks\n")
		}
		if r.AutopilotSettings.RemindStart {
			out += md.Markdownf("    - Remind task users **%v** prior to start\n", r.AutopilotSettings.RemindStartPrior)
		}
		if r.AutopilotSettings.RemindFinish {
			out += md.Markdownf("    - Remind task users **%v** prior to finish\n", r.AutopilotSettings.RemindFinishPrior)
		}
	} else {
		out += md.Markdownf("  - Autopilot: **off**\n")
	}

	return out
}

func (r *Rotation) FindUsers(mattermostUserIDs *types.IDSet) []*User {
	uu := []*User{}
	for _, id := range r.MattermostUserIDs.IDs() {
		uu = append(uu, r.Users.Get(id))
	}
	return uu
}

func (r *Rotation) newTicket(defaultID string) *Task {
	r.TaskSettings.Seq++
	t := NewTask(r.RotationID)
	if defaultID == "" {
		defaultID = strconv.Itoa(r.TaskSettings.Seq)
	}
	t.TaskID = types.ID(r.Name() + "#" + defaultID)
	t.Require = r.TaskSettings.Require.Clone()
	t.Limit = r.TaskSettings.Limit.Clone()
	t.Grace = r.TaskSettings.Grace
	t.ExpectedDuration = r.TaskSettings.Duration
	if t.ExpectedDuration == 0 {
		t.ExpectedDuration = 30 * time.Minute
	}
	return t
}

func (r *Rotation) makeShift(shiftNumber int, now types.Time) (*Task, error) {
	def := r.TaskSettings
	startTime := r.FillSettings.Period.ForNumber(r.FillSettings.Beginning, shiftNumber)
	nextTime := r.FillSettings.Period.ForNumber(startTime, 1)

	// Check if an overlapping shift already exists
	int := types.NewInterval(startTime, nextTime)
	for _, t := range r.Tasks.AsArray() {
		tInt := types.NewDurationInterval(t.ExpectedStart, t.ExpectedDuration)
		switch t.State {
		case TaskStateFinished:
			tInt = types.NewInterval(t.ActualStart, t.ActualFinish)
		case TaskStateStarted:
			tInt = types.NewInterval(t.ActualStart, now)
		}

		if int.Overlaps(tInt) {
			return nil, errors.Errorf(
				"failed to make shift #%v (%v to %v): shift %s (%v - %v) already exists, state %s",
				shiftNumber, startTime, nextTime, t.Markdown(), tInt.Start, tInt.Finish, t.State)
		}
	}

	t := NewTask(r.RotationID)
	t.TaskID = types.ID(fmt.Sprintf("%s#%v", r.Name(), shiftNumber))
	t.Require = def.Require.Clone()
	t.Limit = def.Limit.Clone()
	t.Grace = def.Grace
	t.Description = def.Description

	t.ExpectedStart = startTime
	if def.Duration > 0 {
		t.ExpectedDuration = def.Duration
	} else {
		t.ExpectedDuration = nextTime.Sub(startTime.Time)
	}

	return t, nil
}

func (r *Rotation) queryTasks(finclude func(*Task, types.Time) bool, now types.Time) *Tasks {
	tasks := NewTasks()
	for _, t := range r.Tasks.AsArray() {
		if finclude(t, now) {
			tasks.Set(t)
		}
	}
	return tasks
}

func (r *Rotation) isAutopilotRemindFinish(t *Task, now types.Time) bool {
	if t.State == TaskStateStarted && !t.AutopilotRemindedFinish {
		remindTime := t.ActualStart.Add(t.ExpectedDuration - r.AutopilotSettings.RemindFinishPrior)
		if !now.Before(remindTime) {
			return true
		}
	}
	return false
}

func (r *Rotation) isAutopilotRemindStart(t *Task, now types.Time) bool {
	if t.State == TaskStateScheduled && !t.AutopilotRemindedStart {
		remindTime := t.ExpectedStart.Add(-r.AutopilotSettings.RemindStartPrior)
		if !now.Before(remindTime) {
			return true
		}
	}
	return false
}

func (r *Rotation) isAutopilotFinish(t *Task, now types.Time) bool {
	finishTime := t.ActualStart.Add(t.ExpectedDuration)
	return t.State == TaskStateStarted && finishTime.Before(now.Time)
}

func (r *Rotation) isAutopilotStart(t *Task, now types.Time) bool {
	startTime := t.ExpectedStart.Time
	return t.State == TaskStateScheduled && !now.Before(startTime)
}

func (r *Rotation) allTasksForTime(t *Task, now types.Time) bool {
	return !now.Before(t.ExpectedStart.Time) && now.Before(t.ExpectedStart.Add(t.ExpectedDuration))
}

func (r *Rotation) isAutopilotSchedule(t *Task, now types.Time) bool {
	scheduleTime := t.ExpectedStart.Time.Add(-r.AutopilotSettings.SchedulePrior)
	return t.State == TaskStatePending && !now.Before(scheduleTime)
}

func (as AutopilotSettings) isOn() bool {
	return as.Create || as.RemindFinish || as.RemindStart || as.Schedule || as.StartFinish
}
