// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func (c *Command) archiveRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := flag.NewFlagSet("archiveRotation", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	err = c.API.ArchiveRotation(rotation)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", rotation.Name)
	}

	return "Deleted rotation " + rotation.Name, nil
}

func (c *Command) debugDeleteRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := flag.NewFlagSet("debugDeleteRotation", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}

	err = c.API.DebugDeleteRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", rotationID)
	}

	return "Deleted rotation " + rotationID, nil
}
