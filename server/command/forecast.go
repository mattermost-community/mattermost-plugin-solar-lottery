// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) forecast(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandSchedule: c.forecastSchedule,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `forecast schedule|heat`")
}

func (c *Command) forecastSchedule(parameters []string) (string, error) {
	var rotationID, rotationName string
	autofill := false
	startDate := time.Now().Format(api.DateFormat)
	numShifts := 3
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&autofill, flagAutofill, false, "automatically fill shifts that are not scheduled yet")
	fs.StringVar(&startDate, flagStart, startDate, "starting with date, with --schedule")
	fs.IntVar(&numShifts, flagShifts, numShifts, "number of shifts to scan, with --schedule")
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("forecast schedule", fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	shifts, err := c.API.Forecast(rotation, startDate, numShifts, autofill)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("Rotation %s %v shifts, starting %v:\n", api.MarkdownRotation(rotation), numShifts, startDate)
	for shiftNumber, shift := range shifts {
		if shift != nil {
			out += "- " + api.MarkdownShift(shiftNumber, shift) + "\n"
		}
	}
	return out, nil
}
