// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type TaskService interface {
	MakeTicket(rotationID types.ID, summary, description string) (*Task, error)
}

func (sl *sl) MakeTicket(rotationID types.ID, summary, description string) (*Task, error) {
	err := sl.Setup(pushLogger("MakeTicket", bot.LogContext{}))
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	var t *Task
	_, err = sl.UpdateRotation(rotationID, func(r *Rotation) error {
		t = r.TaskMaker.NewTicket(rotationID, "")
		t.Summary = summary
		t.Description = description
		id, err := sl.Store.Entity(KeyTask).NewID(string(t.TaskID))
		if err != nil {
			return err
		}
		t.TaskID = id

		err = sl.storeTask(t)
		if err != nil {
			return err
		}

		r.PendingTaskIDs.Set(t.TaskID)
		if r.expandedTasks {
			r.pending.Set(t)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (sl *sl) loadTasks(taskIDs *types.IDSet) (*Tasks, error) {
	tasks := NewTasks()
	for _, id := range taskIDs.IDs() {
		t, err := sl.loadTask(id)
		if err != nil {
			return nil, err
		}
		tasks.Set(t)
	}
	return tasks, nil
}

func (sl *sl) loadTask(taskID types.ID) (*Task, error) {
	t := NewTask("")
	err := sl.Store.Entity(KeyTask).Load(t.TaskID, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (sl *sl) storeTask(t *Task) error {
	t.PluginVersion = sl.conf.PluginVersion
	return sl.Store.Entity(KeyTask).Store(t.TaskID, t)
}
