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

type TaskState string

var ErrWrongState = errors.New("operation is not allowed in this state")

const (
	TaskStatePending    = TaskState("pending")
	TaskStateScheduled  = TaskState("scheduled")
	TaskStateInProgress = TaskState("inprogress")
	TaskStateFinished   = TaskState("finished")
)

type Task struct {
	//TODO set PluginID on save
	PluginVersion string
	TaskID        types.ID
	RotationID    types.ID
	State         TaskState
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

func (t Task) Markdown() md.MD {
	return md.Markdownf("%s", t.String())
}

func (t Task) String() string {
	return fmt.Sprintf("%s", t.TaskID)
}

func (t *Task) NewUnavailable() []*Unavailable {
	interval := t.Actual
	if interval.IsEmpty() {
		interval = t.Scheduled
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
			TaskID:     t.TaskID,
			RotationID: t.RotationID,
		},
	}

	if t.Grace > 0 {
		uu = append(uu, &Unavailable{
			Reason: ReasonGrace,
			Interval: types.Interval{
				Start:  interval.Finish,
				Finish: types.NewTime(interval.Finish.Add(t.Grace)),
			},
			TaskID: t.TaskID,
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
