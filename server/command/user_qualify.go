// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
)

func (c *Command) qualifyUsers(parameters []string) (string, error) {
	var skillName string
	var level sl.Level
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withSkillFlags(fs, &skillName, &level)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	if skillName == "" || level == 0 {
		return c.flagUsage(fs), errors.New("must provide --level and --skill values")
	}

	users, err := c.users(fs.Args())
	if err != nil {
		return "", err
	}

	err = c.SL.Qualify(users, skillName, level)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Qualified %s as %s", users.Markdown(), sl.MarkdownSkillLevel(skillName, level)), nil
}
