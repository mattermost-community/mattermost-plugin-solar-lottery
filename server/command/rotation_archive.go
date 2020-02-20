// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
)

func (c *Command) archiveRotation(parameters []string) (string, error) {
	fs := newRotationUsersFlagSet()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs.FlagSet), err
	}

	r, _, err := c.rotationUsers(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.ArchiveRotation(r)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", r.Name())
	}

	return "Archived rotation " + r.Name(), nil
}

func (c *Command) debugDeleteRotation(parameters []string) (string, error) {
	fs := newRotationUsersFlagSet()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs.FlagSet), err
	}

	r, _, err := c.rotationUsers(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.DebugDeleteRotation(r.RotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", r.Markdown())
	}

	return "Deleted rotation " + r.Markdown(), nil
}
