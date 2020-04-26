// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) userLeave(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, mattermostUserIDs, err := c.resolveRotationUsernames()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.LeaveRotation(sl.InJoinRotation{
			MattermostUserIDs: mattermostUserIDs,
			RotationID:        rotationID,
		}))
}
