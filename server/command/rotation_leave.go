// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c *Command) leaveRotation(parameters []string) (string, error) {
	fs := newRotationUsersFlagSet()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs.FlagSet), err
	}

	r, users, err := c.rotationUsers(fs)

	deleted, err := c.SL.LeaveRotation(r, users)
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", deleted.Markdown())
	}

	return fmt.Sprintf("%s removed from rotation %s", deleted.Markdown(), r.Markdown()), nil
}
