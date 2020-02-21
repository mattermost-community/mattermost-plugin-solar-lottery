// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) disqualifyUsers(parameters []string) (string, error) {
	fs := newFS()
	jsonOut := fJSON(fs)
	skill := fSkill(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	users, err := c.loadUsernames(fs.Args())
	if err != nil {
		return "", err
	}

	err = c.SL.Disqualify(users, *skill)
	if err != nil {
		return "", err
	}
	if *jsonOut {
		return md.JSONBlock(users), nil
	}
	return fmt.Sprintf("Disqualified %s from %s", users.Markdown(), *skill), nil
}
