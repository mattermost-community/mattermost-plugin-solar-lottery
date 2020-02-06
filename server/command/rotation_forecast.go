// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) forecastRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	start, numShifts, sampleSize := 0, 3, 10
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.IntVarP(&numShifts, flagNumber, flagPNumber, numShifts, "number of shifts to forecast")
	fs.IntVarP(&start, flagStart, flagPStart, start, "number of shifts to forecast")
	fs.IntVar(&sampleSize, flagSampleSize, sampleSize, "number of guesses to run")
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

	forecast, err := c.SL.ForecastRotation(rotation, start, numShifts, sampleSize)
	if err != nil {
		return "", err
	}

	return utils.JSONBlock(forecast), nil
}
