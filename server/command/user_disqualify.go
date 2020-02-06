// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"
)

func (c *Command) disqualifyUsers(parameters []string) (string, error) {
	var usernames, skillName string
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to disqualify from skill")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	err = c.SL.Disqualify(usernames, skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Disqualified %s from %s", usernames, skillName), nil
}
