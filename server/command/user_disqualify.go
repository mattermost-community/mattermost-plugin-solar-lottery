// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) disqualifyUsers(parameters []string) (string, error) {
	fs := newFS()
	jsonOut := fJSON(fs)
	skill := fSkill(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	mattermostUserIDs, err := c.resolveUsernames(fs.Args())
	if err != nil {
		return "", err
	}
	disqualified, err := c.SL.Disqualify(mattermostUserIDs, types.ID(*skill))
	if err != nil {
		return "", err
	}
	if *jsonOut {
		return md.JSONBlock(disqualified), nil
	}
	return fmt.Sprintf("Disqualified %s from %s", disqualified.Markdown(), *skill), nil
}
