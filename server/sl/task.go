// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

var ErrWrongState = errors.New("operation is not allowed in this state")

const (
	// New tasks that have been submitted are Pending. There are no restrictions
	// on assigning, filling, or un-assigning users to pending tasks.
	// No DMs are sent to users added to or removed from these tasks.
	TaskStatePending = types.ID("pending")

	// Scheduled tasks are normally verified to have met the requirements and
	// constraints. Assigned users receive a message when a task is scheduled. A
	// scheduled task may no longer be filled, but users can still be assigned
	// and unassigned manually, with DM notifications. Users added to a
	// scheduled task get task-related unavailability events added to their
	// calendars.
	TaskStateScheduled = types.ID("scheduled")

	// An in-progress task. Users can be assigned, but not filled (nor unassigned?).
	TaskStateStarted = types.ID("started")

	// Finished tasks are archived, and are not yet used.
	TaskStateFinished = types.ID("finished")
)

type Task struct {
	//TODO set PluginID on save
	PluginVersion string
	TaskID        types.ID
	RotationID    types.ID
	State         types.ID
	Created       types.Time
	Summary       string
	Description   string

	Scheduled         types.Interval `json:",omitempty"`
	Require           *Needs         `json:",omitempty"`
	Limit             *Needs         `json:",omitempty"`
	Actual            types.Interval `json:",omitempty"`
	Grace             time.Duration  `json:",omitempty"`
	MattermostUserIDs *types.IDSet   `json:",omitempty"`

	Users *Users `json:"-"`
}

func NewTask(rotationID types.ID) *Task {
	return &Task{
		State:             TaskStatePending,
		Created:           types.NewTime(),
		RotationID:        rotationID,
		Require:           NewNeeds(),
		Limit:             NewNeeds(),
		MattermostUserIDs: types.NewIDSet(),
		Users:             NewUsers(),
	}
}

func (t Task) GetID() types.ID {
	return t.TaskID
}

func (t Task) MarkdownBullets(rotation *Rotation) md.MD {
	out := md.Markdownf("- %s\n", t.Markdown())
	out += md.Markdownf("  - Status: **%s**\n", t.State)
	out += md.Markdownf("  - Users: **%v**\n", t.MattermostUserIDs.Len())
	// for _, user := range rotation.TaskUsers(&t) {
	// 	out += fmt.Sprintf("    - %s\n", user.MarkdownWithSkills())
	// }
	return out
}

func (task Task) Markdown() md.MD {
	return md.Markdownf("%s", task.String())
}

func (task Task) String() string {
	return fmt.Sprintf("%s", task.TaskID)
}

func (task *Task) NewUnavailable() []*Unavailable {
	interval := task.Actual
	if interval.IsEmpty() {
		interval = task.Scheduled
	}
	if interval.IsEmpty() {
		now := types.NewTime()
		interval = types.Interval{
			Start:  now,
			Finish: now,
		}
	}
	uu := []*Unavailable{
		{
			Reason:     ReasonTask,
			Interval:   interval,
			TaskID:     task.TaskID,
			RotationID: task.RotationID,
		},
	}

	if task.Grace > 0 {
		uu = append(uu, &Unavailable{
			Reason: ReasonGrace,
			Interval: types.Interval{
				Start:  interval.Finish,
				Finish: types.NewTime(interval.Finish.Add(task.Grace)),
			},
			TaskID: task.TaskID,
		})
	}

	return uu
}

func (task *Task) isReadyToStart() (ready bool, whyNot string, err error) {
	if task.State != TaskStatePending && task.State != TaskStateScheduled {
		return false, "", errors.Wrap(ErrWrongState, string(task.State))
	}

	unmetNeeds := task.Require.Unmet(task.Users)
	if unmetNeeds.IsEmpty() {
		return true, "", nil
	}

	whyNot = FillError{
		UnmetNeeds: unmetNeeds,
		Err:        errors.New("not filled"),
		TaskID:     task.TaskID,
	}.Error()

	return false, whyNot, nil
}
