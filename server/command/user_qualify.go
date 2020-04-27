// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) userQualify(parameters []string) (md.MD, error) {
	skills := c.flags().StringSliceP("skills", "s", nil, "skills, with optional levels (1-4) as in `--skills=web-3,server-2`.")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	var skillLevels []sl.SkillLevel
	for _, s := range *skills {
		skillLevel := sl.SkillLevel{}
		err = skillLevel.Set(s)
		if err != nil {
			return c.flagUsage(), err
		}
		skillLevels = append(skillLevels, skillLevel)
	}

	mattermostUserIDs, err := c.resolveUsernames(c.flags().Args())
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.Qualify(sl.InQualify{
			MattermostUserIDs: mattermostUserIDs,
			SkillLevels:       skillLevels,
		}))
}
