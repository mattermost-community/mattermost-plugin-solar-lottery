// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package queue

import (
	"github.com/pkg/errors"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const Type = "queue"

type taskFiller struct{}

var _ sl.TaskFiller = (*taskFiller)(nil)

func New() sl.TaskFiller {
	return &taskFiller{}
}

func (*taskFiller) FillTask(r *sl.Rotation, task *sl.Task, now types.Time, logger bot.Logger) (*sl.Users, error) {
	return nil, errors.New("Queue autofill is not implemented")
}
