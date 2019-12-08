// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/pkg/errors"
)

func (c *Command) skill(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}

	var skills []string
	var err error
	switch parameters[0] {
	case "add":
		if len(parameters) < 2 {
			return "", errors.New("invalid syntax TODO")
		}
		skills, err = c.API.AddSkill(parameters[1])

	case "delete":
		if len(parameters) < 2 {
			return "", errors.New("invalid syntax TODO")
		}
		skills, err = c.API.DeleteSkill(parameters[1])

	case "list":
		skills, err = c.API.ListSkills()
	}
	if err != nil {
		return "", err
	}

	return "Known skills: " + utils.JSONBlock(skills), nil
}
