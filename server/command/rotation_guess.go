// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
)

func (c *Command) guessRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	start, numShifts := 0, 3
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.IntVarP(&numShifts, flagNumber, flagPNumber, numShifts, "number of shifts to forecast")
	fs.IntVarP(&start, flagStart, flagPStart, start, "number of shifts to forecast")
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

	shifts, err := c.SL.Guess(rotation, start, numShifts)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("Rotation %s %v shifts, starting %v:\n", rotation.Markdown(), numShifts, start)
	for _, shift := range shifts {
		if shift != nil {
			out += shift.MarkdownBullets(rotation)
		}
	}
	return out, nil
}
