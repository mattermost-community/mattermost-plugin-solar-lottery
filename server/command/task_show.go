// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) taskShow(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	taskID, _, err := c.resolveTaskIDUsernames()
	if err != nil {
		return "", err
	}

	task, err := c.SL.LoadTask(taskID)
	return c.normalOut(task, err)
}
