// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func (c *Command) user(parameters ...string) (string, error) {
	subcommands := map[string]func(...string) (string, error){
		"show":  c.showUser,
		"skill": c.skillUser,
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

func (c *Command) skillUser(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	skill := parameters[0]

	s := flag.NewFlagSet("skill", flag.ContinueOnError)
	level := ""
	deleteSkill := false
	s.StringVar(&level, "level", "", "beginner, intermediate, advanced, or expert")
	s.BoolVar(&deleteSkill, "delete", false, "specify to remove the skill")
	err := s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}
	if deleteSkill && level != "" {
		return "", errors.New("`--delete` is not compatible with other flags")
	}
	if !deleteSkill && level == "" {
		return "", errors.New("must provide `--level`")
	}

	var user *store.User
	if deleteSkill {
		user, err = c.API.DeleteUserSkill(skill)
	} else {
		var l api.Level
		l, err = api.ParseLevel(level)
		if err != nil {
			return "", err
		}
		user, err = c.API.UpdateUserSkill(skill, int(l))
	}
	if err != nil {
		return "", err
	}

	return utils.JSONBlock(user), nil
}
