// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) qualifyUsers(parameters []string) (string, error) {
	fs := newFS()
	jsonOut := fJSON(fs)
	skillLevel := fSkillLevel(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	if skillLevel.Skill == "" || skillLevel.Level == 0 {
		return c.flagUsage(fs), errors.New("must provide --level and --skill values")
	}

	users, err := c.loadUsernames(fs.Args())
	if err != nil {
		return "", err
	}

	err = c.SL.Qualify(users, *skillLevel)
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(users), nil
	}
	return fmt.Sprintf("Qualified %s as %s", users.Markdown(), skillLevel.String()), nil
}
