// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/spf13/pflag"
)

func (c *Command) qualifyUsers(parameters []string) (string, error) {
	var usernames, skillName string
	var level api.Level
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withSkillFlags(fs, &skillName, &level)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to qualify")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	err = c.API.Qualify(usernames, skillName, level)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Qualified %s as %s", usernames, api.MarkdownSkillLevel(skillName, level)), nil
}
