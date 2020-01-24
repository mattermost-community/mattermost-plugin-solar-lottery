// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
)

func (c *Command) archiveRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := newRotationFlagSet(&rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.SL.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	err = c.SL.ArchiveRotation(rotation)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", rotation.Name)
	}

	return "Deleted rotation " + rotation.Name, nil
}

func (c *Command) debugDeleteRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := newRotationFlagSet(&rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}

	err = c.SL.DebugDeleteRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", rotationID)
	}

	return "Deleted rotation " + rotationID, nil
}
