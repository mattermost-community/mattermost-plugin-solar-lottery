// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) taskUnassign(parameters []string) (md.MD, error) {
	force := c.flags().BoolP("force", "f", false, "ignore task status constraints")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	taskID, mattermostUserIDs, err := c.resolveTaskIDUsernames()
	if err != nil {
		return "", err
	}

	return c.normalOut(c.SL.UnassignTask(sl.InAssignTask{
		TaskID:            taskID,
		MattermostUserIDs: mattermostUserIDs,
		Force:             *force,
	}))
}
