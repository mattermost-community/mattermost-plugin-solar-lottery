// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) showUser(parameters []string) (string, error) {
	fs := newFS()
	_ = fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	users, err := c.loadUsernames(fs.Args())
	if err != nil {
		return "", err
	}

	return md.JSONBlock(users), nil
}
