// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/spf13/pflag"
)

func (c *Command) showUser(parameters []string) (string, error) {
	usernames := ""
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	users, err := c.API.LoadMattermostUsers(usernames)
	if err != nil {
		return "", err
	}
	return utils.JSONBlock(users), nil
}
