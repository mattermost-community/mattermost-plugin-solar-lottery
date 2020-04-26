// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) userDisqualify(parameters []string) (md.MD, error) {
	skills := c.flags().StringSliceP("skills", "s", nil, "skills to disqualify from, e.g. `--skills=web,server`")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	mattermostUserIDs, err := c.resolveUsernames(c.flags().Args())
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.Disqualify(sl.InDisqualify{
			MattermostUserIDs: mattermostUserIDs,
			Skills:            *skills,
		}))
}
