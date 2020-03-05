// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

type TaskService interface {
	MakeTicket(rotationID types.ID, summary, description string) (*Task, error)
	AssignTask(types.ID, *types.IDSet, bool) (*Task, *Users, error)
	FillTask(*Task) error
}

func (sl *sl) MakeTicket(rotationID types.ID, summary, description string) (*Task, error) {
	err := sl.Setup(pushLogger("MakeTicket", bot.LogContext{}))
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	var t *Task
	_, err = sl.UpdateRotation(rotationID, func(r *Rotation) error {
		t = r.TaskMaker.newTicket(r, "")
		t.Summary = summary
		t.Description = description
		var id types.ID
		id, err = sl.Store.Entity(KeyTask).NewID(string(t.TaskID))
		if err != nil {
			return err
		}
		t.TaskID = id

		err = sl.storeTask(t)
		if err != nil {
			return err
		}
		r.TaskIDs.Set(t.TaskID)
		if r.pending != nil {
			r.pending.Set(t)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sl.Infof("%s created ticket %s.", sl.actingUser.Markdown(), t.Markdown())
	return t, nil
}

func (sl *sl) AssignTask(taskID types.ID, mattermostUserIDs *types.IDSet, force bool) (*Task, *Users, error) {
	users := NewUsers()
	task := NewTask("")
	err := sl.Setup(
		pushLogger("AssignTask", bot.LogContext{ctxForce: force}),
		withLoadTask(taskID, &task),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return nil, nil, err
	}
	defer sl.popLogger()

	max := task.Max
	added := NewUsers()
	for _, user := range users.AsArray() {
		if task.MattermostUserIDs.Contains(user.MattermostUserID) {
			continue
		}
		if !force {
			var failed *Needs
			max, failed = user.checkMaxConstraints(max)
			if !failed.IsEmpty() {
				return nil, nil, errors.Errorf("user %s failed max constraints %s", user.Markdown(), failed.MarkdownSkillLevels())
			}
		}
		task.MattermostUserIDs.Set(user.MattermostUserID)
		if task.users != nil {
			task.users.Set(user)
		}
		added.Set(user)
	}
	err = sl.storeTask(task)
	if err != nil {
		return nil, nil, err
	}

	sl.Infof("%s assigned %s to ticket %s.",
		sl.actingUser.Markdown(), added.Markdown(), task.Markdown())
	return task, added, nil
}

func (sl *sl) FillTask(*Task) error {
	return errors.New("<><> TODO")
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
	err := sl.Store.Entity(KeyTask).Load(taskID, t)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load %s", taskID)
	}
	return t, nil
}

func (sl *sl) storeTask(t *Task) error {
	t.PluginVersion = sl.conf.PluginVersion
	err := sl.Store.Entity(KeyTask).Store(t.TaskID, t)
	if err != nil {
		return errors.Wrapf(err, "failed to store %s", t.String())
	}
	return nil
}

func (sl *sl) expandTaskUsers(t *Task) error {
	if t.users != nil {
		return nil
	}
	users, err := sl.loadStoredUsers(t.MattermostUserIDs)
	if err != nil {
		return err
	}
	err = sl.expandUsers(users)
	if err != nil {
		return err
	}
	t.users = users
	return nil
}
