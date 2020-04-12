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
		withExpandedTask(&params.TaskID, task),
		withExpandedRotation(&task.RotationID, r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	filled, err := sl.fillTask(r, task, params.Time)
	if err != nil {
		return nil, err
	}

	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	lastServed := task.ExpectedStart
	if lastServed.IsZero() {
		lastServed = params.Time
	}
	for _, user := range filled.AsArray() {
		user.LastServed.Set(r.RotationID, lastServed.Unix())
	}
	err = sl.storeUsers(filled)
	if err != nil {
		return nil, err
	}

	out := &OutAssignTask{
		MD:      md.Markdownf("Auto-assigned %s to ticket %s.", filled.Markdown(), task.Markdown()),
		Task:    task,
		Changed: filled,
	}
	sl.logAPI(out)
	return out, nil
}
