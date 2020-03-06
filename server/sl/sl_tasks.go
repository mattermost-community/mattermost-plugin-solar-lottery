// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

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

func (sl *sl) assignTask(task *Task, users *Users, force bool) (*Users, error) {
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
				return nil, errors.Errorf("user %s failed max constraints %s", user.Markdown(), failed.MarkdownSkillLevels())
			}
		}
		task.MattermostUserIDs.Set(user.MattermostUserID)
		if task.users != nil {
			task.users.Set(user)
		}
		added.Set(user)
	}
	return added, nil
}

func (sl *sl) fillTask(r *Rotation, task *Task) (added *Users, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "failed to fill task %s", task.Markdown())
		}
	}()

	// Autofill is only allowed on pending tasks
	if task.State != TaskStatePending {
		return nil, errors.Wrap(ErrWrongState, string(task.State))
	}

	filler, err := sl.taskFiller(r)
	if err != nil {
		return nil, err
	}
	added, err = filler.FillTask(r, task, sl.Logger)
	if err != nil {
		return nil, err
	}
	return added, nil
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
