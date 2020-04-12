// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InAssignTask struct {
	TaskID            types.ID
	MattermostUserIDs *types.IDSet
	Force             bool
	Time              types.Time
}

type OutAssignTask struct {
	md.MD
	Task    *Task
	Changed *Users
}

func (sl *sl) AssignTask(params InAssignTask) (*OutAssignTask, error) {
	users := NewUsers()
	task := NewTask("")
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("AssignTask", params),
		withExpandedTask(&params.TaskID, task),
		withExpandedRotation(&task.RotationID, r),
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

	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	out := &OutAssignTask{
		MD:      md.Markdownf("assigned %s to ticket %s.", assigned.Markdown(), task.Markdown()),
		Task:    task,
		Changed: assigned,
	}
	sl.logAPI(out)
	return out, nil
}
