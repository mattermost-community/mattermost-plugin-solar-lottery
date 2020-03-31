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
	PluginVersion  string
	RotationID     types.ID
	IsArchived     bool
	TaskFillerType types.ID

	Type        types.ID     // ticket or shift
	TaskSeq     int          `json:",omitempty"`
	Beginning   types.Time   `json:",omitempty"`
	ShiftPeriod types.Period `json:",omitempty"`

	// Task defaults
	Require     *Needs        `json:",omitempty"`
	Limit       *Needs        `json:",omitempty"`
	Duration    time.Duration `json:",omitempty"`
	Grace       time.Duration `json:",omitempty"`
	Description string        `json:",omitempty"`

	MattermostUserIDs *types.IDSet `json:",omitempty"`
	TaskIDs           *types.IDSet `json:",omitempty"`

	loaded bool
	Users  *Users `json:"-"`
	tasks  *Tasks
}

const (
	TypeTicket = types.ID("Ticket")
	TypeShift  = types.ID("Shift")
)

func NewRotation() *Rotation {
	r := &Rotation{}
	r.init()
	return r
}

func (r *Rotation) init() {
	if r.MattermostUserIDs == nil {
		r.MattermostUserIDs = types.NewIDSet()
	}
	if r.TaskIDs == nil {
		r.TaskIDs = types.NewIDSet()
	}
	if r.Require == nil {
		r.Require = NewNeeds()
	}
	if r.Limit == nil {
		r.Limit = NewNeeds()
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

	out += md.Markdownf("  - Type: **%s**\n", r.Type)
	out += md.Markdownf("  - Require: **%s**\n", r.Require.Markdown())
	out += md.Markdownf("  - Limit: **%v**\n", r.Limit.Markdown())
	out += md.Markdownf("  - Grace: %v\n", r.Grace)
	if r.Type == TypeShift {
		out += md.Markdownf("  - Shift period: **%v**\n", r.ShiftPeriod.String())
		out += md.Markdownf("  - Shifts beginning: **%v**\n", r.Beginning)
	}

	// if rotation.Autopilot.On {
	// 	out += fmt.Sprintf("  - Autopilot: **on**\n")
	// 	out += fmt.Sprintf("    - Auto-start: **%v**\n", rotation.Autopilot.StartFinish)
	// 	out += fmt.Sprintf("    - Auto-fill: **%v**, %v days prior to start\n", rotation.Autopilot.Fill, rotation.Autopilot.FillPrior)
	// 	out += fmt.Sprintf("    - Notify users in advance: **%v**, %v days prior to transition\n", rotation.Autopilot.Notify, rotation.Autopilot.NotifyPrior)
	// } else {
	// 	out += fmt.Sprintf("  - Autopilot: **off**\n")
	// }

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
	r.TaskSeq++
	t := NewTask(r.RotationID)
	if defaultID == "" {
		defaultID = strconv.Itoa(r.TaskSeq)
	}
	t.TaskID = types.ID(r.Name() + "#" + defaultID)
	t.Require = r.Require.Clone()
	t.Limit = r.Limit.Clone()
	t.Grace = r.Grace
	t.ExpectedDuration = r.Duration
	return t
}

func (r *Rotation) makeShift(shiftNumber int) (*Task, error) {
	startTime := r.ShiftPeriod.Next(r.Beginning, shiftNumber-1)
	nextTime := r.ShiftPeriod.Next(startTime, 1)

	// Check if an overlapping shift already exists
	int := types.NewInterval(startTime, nextTime)
	for _, t := range r.tasks.AsArray() {
		tInt := types.NewDurationInterval(t.ExpectedStart, t.ExpectedDuration)
		switch t.State {
		case TaskStateFinished:
			tInt = types.NewInterval(t.ActualStart, t.ActualFinish)
		case TaskStateStarted:
			tInt = types.NewInterval(t.ActualStart, types.NewTime(time.Now()))
		}

		if int.Overlaps(tInt) {
			return nil, errors.Errorf(
				"failed to make shift #%v (%v to %v): shift %s (%v - %v) already exists, state %s",
				shiftNumber, startTime, nextTime, t.Markdown(), tInt.Start, tInt.Finish, t.State)
		}
	}

	t := NewTask(r.RotationID)
	t.TaskID = types.ID(fmt.Sprintf("%s#%v", r.Name(), shiftNumber))
	t.Require = r.Require.Clone()
	t.Limit = r.Limit.Clone()
	t.Grace = r.Grace
	t.Description = r.Description

	t.ExpectedStart = startTime
	if r.Duration > 0 {
		t.ExpectedDuration = r.Duration
	} else {
		t.ExpectedDuration = nextTime.Sub(startTime.Time)
	}

	return t, nil
}
