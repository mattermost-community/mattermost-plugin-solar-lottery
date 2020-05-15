// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) userJoin(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	starting, err := c.withTimeFlag("starting", fmt.Sprintf("time for user to start participating"))
	if err != nil {
		return c.flagUsage(), err
	}
	err = c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, mattermostUserIDs, err := c.resolveRotationUsernames()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.JoinRotation(sl.InJoinRotation{
			MattermostUserIDs: mattermostUserIDs,
			RotationID:        rotationID,
			Starting:          *starting,
		}))
}
