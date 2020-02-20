// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"
)

func (c *Command) disqualifyUsers(parameters []string) (string, error) {
	var skillName string
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	users, err := c.users(fs.Args())
	if err != nil {
		return "", err
	}

	err = c.SL.Disqualify(users, skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Disqualified %s from %s", users.Markdown(), skillName), nil
}
