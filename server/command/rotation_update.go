// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func withRotationUpdateFlags(fs *pflag.FlagSet, size *int, grace *int) {
	fs.IntVar(size, flagSize, 0, "target number of people in each shift. 0 (default) means unlimited, based on needs")
	fs.IntVar(grace, flagGrace, 1, "blocks for serving this many shifts after one served")
}

func (c *Command) updateRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	var size, grace int
	fs := pflag.NewFlagSet("updateRotation", pflag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	withRotationUpdateFlags(fs, &size, &grace)
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

	err = c.API.UpdateRotation(rotation, func(rotation *api.Rotation) error {
		if grace != 0 {
			rotation.Grace = grace
		}
		if size != 0 {
			rotation.Size = size
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return "Updated rotation " + api.MarkdownRotationWithDetails(rotation), nil
}
