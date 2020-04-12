// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InTransitionTask struct {
	TaskID types.ID
	State  types.ID
	Time   types.Time
}

type OutTransitionTask struct {
	md.MD
	Task      *Task
	PrevState types.ID
}

func (sl *sl) TransitionTask(params InTransitionTask) (*OutTransitionTask, error) {
	task := NewTask("")
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("TransitionTask", params),
		withExpandedTask(&params.TaskID, task),
		withExpandedRotation(&task.RotationID, r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	prevState := task.State
	err = sl.transitionTask(r, task, params.Time, params.State)
	if err != nil {
		return nil, err
	}

	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	out := &OutTransitionTask{
		MD:        md.Markdownf("transitioned %s to %s.", task.Markdown(), task.State),
		Task:      task,
		PrevState: prevState,
	}
	sl.logAPI(out)
	return out, nil
}
