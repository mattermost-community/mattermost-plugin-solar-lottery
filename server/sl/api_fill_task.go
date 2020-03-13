// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (sl *sl) FillTask(params InAssignTask) (*OutAssignTask, error) {
	task := NewTask("")
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("FillTask", params),
		withLoadExpandTask(&params.TaskID, task),
		withLoadExpandRotation(&task.RotationID, r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	filled, err := sl.fillTask(r, task)
	if err != nil {
		return nil, err
	}

	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	out := &OutAssignTask{
		MD:    md.Markdownf("Auto-assigned %s to ticket %s.", filled.Markdown(), task.Markdown()),
		Task:  task,
		Added: filled,
	}
	sl.LogAPI(out)
	return out, nil
}
