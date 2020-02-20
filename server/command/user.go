// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/pkg/errors"
)

func (c *Command) user(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandDisqualify:  c.disqualifyUsers,
		commandQualify:     c.qualifyUsers,
		commandShow:        c.showUser,
		commandUnavailable: c.userUnavailable,
		// commandForecast:    c.userForecast,
	}
	return c.handleCommand(subcommands, parameters)
}

func (c *Command) users(args []string) (users sl.UserMap, err error) {
	users = sl.UserMap{}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "@") {
			return nil, errors.New("`@username`'s expected")
		}
		arg = arg[1:]
		user, err := c.SL.LoadMattermostUsername(arg)
		if err != nil {
			return nil, err
		}
		users[user.MattermostUserID] = user
	}

	return users, nil
}
