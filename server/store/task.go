// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

type TaskStatus string

const (
	TaskStatusOpen     = TaskStatus("open")
	TaskStatusStarted  = TaskStatus("started")
	TaskStatusFinished = TaskStatus("finished")
)

type Task struct {
	PluginVersion string
	TaskID        string
	Status        TaskStatus
	Created       utils.Time

	Scheduled *utils.Interval `json:",omitempty"`
	Requires  Needs           `json:",omitempty"`
	Limits    Needs           `json:",omitempty"`

	Actual            *utils.Interval `json:",omitempty"`
	MattermostUserIDs IDMap           `json:",omitempty"`
	Autopilot         TaskAutopilot   `json:",omitempty"`
}

type TaskAutopilot struct {
	Filled         utils.Time `json:",omitempty"`
	NotifiedStart  utils.Time `json:",omitempty"`
	NotifiedFinish utils.Time `json:",omitempty"`
}

func NewTask() *Task {
	return &Task{
		Status:            TaskStatusOpen,
		Created:           utils.TimeNow(),
		MattermostUserIDs: IDMap{},
	}
}
