// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) joinRotation(parameters []string) (string, error) {
	var starting types.Time
	fs := newRotationFS()
	jsonOut := fJSON(fs)
	fs.Var(&starting, flagStart, fmt.Sprintf("time for user to start participating"))
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	r, users, err := c.rotationUsers(fs)
	if err != nil {
		return "", err
	}

	added, err := c.SL.JoinRotation(r, users, starting)
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", added.Markdown())
	}

	if *jsonOut {
		return md.JSONBlock(added), nil
	}
	return fmt.Sprintf("%s added to rotation %s", added.Markdown(), r.Markdown()), nil
}
