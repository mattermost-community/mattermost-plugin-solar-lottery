// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) showUser(parameters []string) (string, error) {
	usernames := ""
	fs := flag.NewFlagSet("", flag.ContinueOnError)
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
