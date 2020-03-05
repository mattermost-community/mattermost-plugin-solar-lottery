// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) assignTask(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	force := fForce(fs)
	fill := fFill(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	taskID, mattermostUserIDs, err := c.resolveTaskIDUsernames(fs)
	if err != nil {
		return "", err
	}

	t, added, err := c.SL.AssignTask(taskID, mattermostUserIDs, *force)
	if err != nil {
		return "", err
	}

	if *fill {
		err := c.SL.FillTask(t)
		if err != nil {
			return "", err
		}
		t, added, err = c.SL.AssignTask(taskID, t.MattermostUserIDs, false)
		if err != nil {
			return "", err
		}
	}

	if *jsonOut {
		return md.JSONBlock(t), nil
	}
	return fmt.Sprintf("%s assigned to %s", added.Markdown(), t), nil
}
