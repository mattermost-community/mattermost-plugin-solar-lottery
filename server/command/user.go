// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) user(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandDisqualify:  c.disqualifyUsers,
		commandQualify:     c.qualifyUsers,
		commandShow:        c.showUser,
		commandUnavailable: c.userUnavailable,
	}
	return c.handleCommand(subcommands, parameters,
		"Usage: `user disqualify|qualify|show|unavailable`. Use `user subcommand --help` for more information.")
}

func (c *Command) showUser(parameters []string) (string, error) {
	usernames := ""
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("user show", fs), err
	}

	users, err := c.API.LoadMattermostUsers(usernames)
	if err != nil {
		return "", err
	}
	return utils.JSONBlock(users), nil
}

func (c *Command) userUnavailable(parameters []string) (string, error) {
	return "TODO", nil
}
