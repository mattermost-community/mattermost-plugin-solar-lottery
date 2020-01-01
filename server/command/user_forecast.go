// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"time"

	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) userForecast(parameters []string) (string, error) {
	var rotationID, rotationName, username string
	numShifts, sampleSize := 12, 10
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.IntVarP(&numShifts, flagNumber, flagPNumber, numShifts, "number of shifts to forecast")
	fs.IntVar(&sampleSize, flagSampleSize, sampleSize, "number of guesses to run")
	fs.StringVarP(&username, flagUsers, flagPUsers, "", "user to forecast (one)")
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

	forecast, err := c.API.ForecastUser(username, rotation, numShifts, sampleSize, time.Now())
	if err != nil {
		return "", err
	}

	return utils.JSONBlock(forecast), nil
}
