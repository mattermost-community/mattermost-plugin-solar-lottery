// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) userShow(parameters []string) (md.MD, error) {
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	mattermostUserIDs, err := c.resolveUsernames(c.flags().Args())
	if err != nil {
		return "", err
	}
	users, err := c.SL.LoadUsers(mattermostUserIDs)
	if err != nil {
		return "", err
	}

	if users.Len() == 1 {
		return md.JSONBlock(users.AsArray()[0]), nil
	}
	return md.JSONBlock(users), nil
}
