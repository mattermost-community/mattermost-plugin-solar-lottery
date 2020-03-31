// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const Type = "solar-lottery"

type taskFiller struct{}

var _ sl.TaskFiller = (*taskFiller)(nil)

func New() sl.TaskFiller {
	return &taskFiller{}
}

func (*taskFiller) FillTask(r *sl.Rotation, task *sl.Task, forTime types.Time, logger bot.Logger) (*sl.Users, error) {
	f := newFill(r, task, forTime, logger)
	return f.fill()
}
