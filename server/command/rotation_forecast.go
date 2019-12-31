// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) forecastRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	start, numShifts, sampleSize := 0, 3, 10
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.IntVarP(&numShifts, flagNumber, flagPNumber, numShifts, "number of shifts to forecast")
	fs.IntVarP(&start, flagStart, flagPStart, start, "number of shifts to forecast")
	fs.IntVar(&sampleSize, flagSampleSize, sampleSize, "number of guesses to run")
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

	forecast, err := c.API.Forecast(rotation, start, numShifts, sampleSize)
	if err != nil {
		return "", err
	}

	return utils.JSONBlock(forecast), nil
}
