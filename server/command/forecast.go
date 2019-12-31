// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) forecast(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandGuess:    c.forecastGuess,
		commandRotation: c.forecastRotation,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `forecast guess|rotation|user`")
}

func (c *Command) forecastGuess(parameters []string) (string, error) {
	var rotationID, rotationName string
	autofill := false
	start, numShifts := 0, 3
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&autofill, flagAutofill, false, "automatically fill shifts that are not scheduled yet")
	fs.IntVarP(&numShifts, flagNumber, flagPNumber, numShifts, "number of shifts to forecast")
	fs.IntVarP(&start, flagStart, flagPStart, start, "number of shifts to forecast")
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
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
			out += "- " + api.MarkdownShift(shiftNumber, shift) + "\n"
		}
	}
	return out, nil
}

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
		return c.subUsage(fs), err
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
