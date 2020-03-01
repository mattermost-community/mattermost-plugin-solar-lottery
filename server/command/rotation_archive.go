// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/pkg/errors"
)

func (c *Command) archiveRotation(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	r, err := c.SL.ArchiveRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", rotationID)
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return "Archived rotation " + r.Markdown(), nil
}

func (c *Command) debugDeleteRotation(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.DebugDeleteRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", rotationID)
	}

	if *jsonOut {
		return md.JSONBlock(rotationID), nil
	}
	return "Deleted rotation " + string(rotationID), nil
}
