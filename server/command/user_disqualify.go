// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) disqualifyUsers(parameters []string) (md.MD, error) {
	skill := c.withFlagSkill()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	mattermostUserIDs, err := c.resolveUsernames(c.fs.Args())
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.Disqualify(sl.InDisqualify{
			MattermostUserIDs: mattermostUserIDs,
			Skill:             types.ID(*skill),
		}))
}
