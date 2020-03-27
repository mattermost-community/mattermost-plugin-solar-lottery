// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) userQualify(parameters []string) (md.MD, error) {
	skillLevel := c.withFlagSkillLevel()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	if skillLevel.Skill == "" || skillLevel.Level == 0 {
		return c.flagUsage(), errors.New("must provide --level and --skill values")
	}

	mattermostUserIDs, err := c.resolveUsernames(c.fs.Args())
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.Qualify(sl.InQualify{
			MattermostUserIDs: mattermostUserIDs,
			SkillLevel:        *skillLevel,
		}))
}
