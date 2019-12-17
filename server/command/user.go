// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) user(parameters ...string) (string, error) {
	users := ""
	fs := flag.NewFlagSet("user", flag.ContinueOnError)
	fs.StringVar(&users, "users", "", "remove other users from rotation.")
	err := fs.Parse(parameters)
	if err != nil {
		return "", err
	}

	switch fs.Arg(0) {
	case "show":
		users, err := c.API.LoadMattermostUsers(users)
		if err != nil {
			return "", err
		}
		return utils.JSONBlock(users), nil

	default:
		return "", errors.Errorf(commandUsage("user show", fs))
	}

}
