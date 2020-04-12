// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (sl *sl) UnassignTask(params InAssignTask) (*OutAssignTask, error) {
	users := NewUsers()
	task := NewTask("")
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("UnassignTask", params),
		withExpandedTask(&params.TaskID, task),
		withExpandedRotation(&task.RotationID, r),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	removed, err := sl.unassignTask(task, users, params.Force)
	if err != nil {
		return nil, err
	}

	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	out := &OutAssignTask{
		MD:      md.Markdownf("assigned %s to ticket %s.", removed.Markdown(), task.Markdown()),
		Task:    task,
		Changed: removed,
	}
	sl.logAPI(out)
	return out, nil
}
