// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/pkg/errors"
)

func (c *Command) rotationArchive(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	r, err := c.SL.ArchiveRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", rotationID)
	}

	return c.normalOut(r, err)
}

func (c *Command) rotationDebugDelete(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	err = c.SL.DebugDeleteRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", rotationID)
	}

	return c.normalOut(md.MD(rotationID), err)
}
