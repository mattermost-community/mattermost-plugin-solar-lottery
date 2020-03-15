// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) joinRotation(parameters []string) (md.MD, error) {
	var starting types.Time
	c.withFlagRotation()
	c.fs.Var(&starting, flagStart, fmt.Sprintf("time for user to start participating"))
	err := c.fs.Parse(parameters)
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
			Starting:          starting,
		}))
}
