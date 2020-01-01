// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) guessRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	autofill := false
	start, numShifts := 0, 3
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.BoolVar(&autofill, flagAutofill, false, "automatically fill shifts that are not scheduled yet")
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
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	shifts, err := c.API.Guess(rotation, start, numShifts, autofill)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("Rotation %s %v shifts, starting %v:\n", api.MarkdownRotation(rotation), numShifts, start)
	for shiftNumber, shift := range shifts {
		if shift != nil {
			out += "- " + api.MarkdownShift(rotation, shiftNumber, shift) + "\n"
		}
	}
	return out, nil
}
