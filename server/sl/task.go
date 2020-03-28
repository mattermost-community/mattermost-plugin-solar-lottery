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

	Start             types.Time     `json:",omitempty"`
	Duration          time.Duration  `json:",omitempty"`
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

func (t Task) Markdown() md.MD {
	return md.Markdownf("%s", t.String())
}

func (t Task) String() string {
	return fmt.Sprintf("%s", t.TaskID)
}

func (t *Task) NewUnavailable() []*Unavailable {
	interval := t.Actual
	if interval.IsEmpty() {
		interval = types.NewDurationInterval(t.Start, t.Duration)
	}
	if interval.IsEmpty() {
		now := types.NewTime()
		interval = types.NewDurationInterval(now, 0)
	}
	uu := []*Unavailable{
		{
			Reason:     ReasonTask,
			Interval:   interval,
			TaskID:     t.TaskID,
			RotationID: t.RotationID,
		},
	}

	if t.Grace > 0 {
		uu = append(uu, &Unavailable{
			Reason:     ReasonGrace,
			Interval:   types.NewDurationInterval(interval.Finish, t.Grace),
			TaskID:     t.TaskID,
			RotationID: t.RotationID,
		})
	}

	return uu
}

func (t *Task) isReadyToStart() (ready bool, whyNot string, err error) {
	if t.State != TaskStatePending && t.State != TaskStateScheduled {
		return false, "", errors.Wrap(ErrWrongState, string(t.State))
	}

	unmetNeeds := t.Require.Unmet(t.Users)
	if unmetNeeds.IsEmpty() {
		return true, "", nil
	}

	whyNot = FillError{
		UnmetNeeds: unmetNeeds,
		Err:        errors.New("not filled"),
		TaskID:     t.TaskID,
	}.Error()

	return false, whyNot, nil
}
