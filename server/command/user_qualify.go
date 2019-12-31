// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) qualifyUsers(parameters []string) (string, error) {
	var usernames, skillName string
	var level api.Level
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	withSkillFlags(fs, &skillName, &level)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
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
