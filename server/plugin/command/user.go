// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/pkg/errors"
)

func (c *Command) user(parameters ...string) (string, error) {
	subcommands := map[string]func(...string) (string, error){
		"show": c.showUser,
	}
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	f := subcommands[parameters[0]]
	if f == nil {
		return "", errors.New("invalid syntax TODO")
	}
	return f(parameters[1:]...)
}

func (c *Command) showUser(parameters ...string) (string, error) {
	user, err := c.API.GetUser()
	if err != nil {
		return "", err
	}
	return utils.JSONBlock(user), nil
}
