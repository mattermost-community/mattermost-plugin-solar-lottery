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

	r, _, err := c.loadRotationUsernames(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.ArchiveRotation(r)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", r.Name())
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return "Archived rotation " + r.Name(), nil
}

func (c *Command) debugDeleteRotation(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	r, _, err := c.loadRotationUsernames(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.DebugDeleteRotation(r.RotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", r.Markdown())
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return "Deleted rotation " + r.Markdown(), nil
}
