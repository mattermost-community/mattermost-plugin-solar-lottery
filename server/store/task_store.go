// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type TaskStore interface {
	LoadTask(taskID string) (*Task, error)
	StoreTask(shift *Task) error
	DeleteTask(taskID string) error
}

func (s *pluginStore) LoadTask(taskID string) (*Task, error) {
	task := NewTask()
	err := kvstore.LoadJSON(s.taskKV, taskID, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *pluginStore) StoreTask(task *Task) error {
	err := kvstore.StoreJSON(s.taskKV, task.TaskID, task)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Task": task,
	}).Debugf("store: Stored task %s", task.TaskID)
	return nil
}

func (s *pluginStore) DeleteTask(taskID string) error {
	err := s.taskKV.Delete(taskID)
	if err != nil {
		return err
	}
	s.Logger.Debugf("store: Deleted task %s %v", taskID)
	return nil
}
