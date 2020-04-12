// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InCreateTicket struct {
	RotationID  types.ID
	Summary     string
	Description string
}

type OutCreateTask struct {
	md.MD
	Task *Task
}

func (sl *sl) CreateTicket(params InCreateTicket) (*OutCreateTask, error) {
	err := sl.Setup(pushAPILogger("CreateTicket", params))
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
		if r.Tasks != nil {
			r.Tasks.Set(task)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := &OutCreateTask{
		MD:   md.Markdownf("created ticket %s.", task.Markdown()),
		Task: task,
	}
	sl.logAPI(out)
	return out, nil
}

type InCreateShift struct {
	RotationID types.ID
	Number     int
}

func (sl *sl) CreateShift(in InCreateShift) (*OutCreateTask, error) {
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("MakeShift", in),
		withLoadRotation(&in.RotationID, r),
		withExpandRotationTasks(r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	t, err := sl.createShift(r, in.Number)
	if err != nil {
		return nil, err
	}
	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return nil, err
	}

	out := &OutCreateTask{
		MD:   md.Markdownf("created shift %s.", t.Markdown()),
		Task: t,
	}
	sl.logAPI(out)
	return out, nil
}
