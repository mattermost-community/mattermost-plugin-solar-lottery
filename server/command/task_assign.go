// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) assignTask(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	force := c.withFlagForce()
	fill := c.withFlagFill()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	taskID, mattermostUserIDs, err := c.resolveTaskIDUsernames()
	if err != nil {
		return "", err
	}

	out, err := c.SL.AssignTask(sl.InAssignTask{
		// return c.normalOut(c.SL.AssignTask(sl.InAssignTask{
		TaskID:            taskID,
		MattermostUserIDs: mattermostUserIDs,
		Fill:              *fill,
		Force:             *force,
	})

	return c.normalOut(out, err)
}
