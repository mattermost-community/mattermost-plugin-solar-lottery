// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InMakeTicket struct {
	RotationID  types.ID
	Summary     string
	Description string
}

type OutMakeTask struct {
	md.MD
	Task *Task
}

func (sl *sl) MakeTicket(params InMakeTicket) (*OutMakeTask, error) {
	err := sl.Setup(pushAPILogger("MakeTicket", params))
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	var task *Task
	_, err = sl.UpdateRotation(params.RotationID, func(r *Rotation) error {
		task = r.newTicket("")
		task.Summary = params.Summary
		task.Description = params.Description
		var id types.ID
		id, err = sl.Store.Entity(KeyTask).NewID(string(task.TaskID))
		if err != nil {
			return err
		}
		task.TaskID = id

		err = sl.storeTask(task)
		if err != nil {
			return err
		}
		r.TaskIDs.Set(task.TaskID)
		if r.tasks != nil {
			r.tasks.Set(task)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := &OutMakeTask{
		MD:   md.Markdownf("created ticket %s.", task.Markdown()),
		Task: task,
	}
	sl.LogAPI(out)
	return out, nil
}

type InMakeShift struct {
	RotationID types.ID
	Number     int
}

func (sl *sl) MakeShift(params InMakeShift) (*OutMakeTask, error) {
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("MakeShift", params),
		withLoadRotation(&params.RotationID, r),
		withExpandRotationTasks(r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	var task *Task
	task, err = r.makeShift(params.Number)
	if err != nil {
		return nil, err
	}
	var id types.ID
	id, err = sl.Store.Entity(KeyTask).NewID(string(task.TaskID))
	if err != nil {
		return nil, err
	}
	task.TaskID = id
	err = sl.storeTask(task)
	if err != nil {
		return nil, err
	}

	r.TaskIDs.Set(task.TaskID)
	if r.tasks != nil {
		r.tasks.Set(task)
	}
	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return nil, err
	}

	out := &OutMakeTask{
		MD:   md.Markdownf("created shift %s.", task.Markdown()),
		Task: task,
	}
	sl.LogAPI(out)
	return out, nil
}
