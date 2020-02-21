// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) leaveRotation(parameters []string) (string, error) {
	fs := newRotationFS()
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	r, users, err := c.rotationUsers(fs)
	if err != nil {
		return "", err
	}

	deleted, err := c.SL.LeaveRotation(r, users)
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", deleted.Markdown())
	}

	if *jsonOut {
		return md.JSONBlock(deleted), nil
	}
	return fmt.Sprintf("%s removed from rotation %s", deleted.Markdown(), r.Markdown()), nil
}
