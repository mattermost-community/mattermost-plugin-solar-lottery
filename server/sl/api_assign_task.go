// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InAssignTask struct {
	TaskID            types.ID
	MattermostUserIDs *types.IDSet
	Fill              bool
	Force             bool
}

type OutAssignTask struct {
	md.MD
	Task     *Task
	Assigned *Users
	Filled   *Users
}

func (sl *sl) AssignTask(params InAssignTask) (*OutAssignTask, error) {
	users := NewUsers()
	task := NewTask("")
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("AssignTask", params),
		withLoadTask(&params.TaskID, task),
		withLoadRotation(&task.RotationID, r),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	assigned, err := sl.assignTask(task, users, params.Force)
	if err != nil {
		return nil, err
	}

	filled := NewUsers()
	if params.Fill {
		filled, err = sl.fillTask(r, task)
		if err != nil {
			return nil, err
		}
	}

	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	txt, sep := "", ""
	if !assigned.IsEmpty() {
		txt += fmt.Sprintf("assigned %s", assigned.Markdown())
		sep = ", "
	}
	if !filled.IsEmpty() {
		txt += sep + fmt.Sprintf("auto-filled %s", filled.Markdown())
	}
	txt += fmt.Sprintf(" to ticket %s.", task.Markdown())
	out := &OutAssignTask{
		MD:       md.MD(txt),
		Task:     task,
		Assigned: assigned,
		Filled:   filled,
	}

	sl.LogAPI(out)
	return out, nil
}
