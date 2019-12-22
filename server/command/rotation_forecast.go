// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) forecast(parameters ...string) (string, error) {
	fs := flag.NewFlagSet("forecast", flag.ContinueOnError)
	autofill := false
	startDate := time.Now().Format(api.DateFormat)
	numShifts := 3
	fs.BoolVar(&autofill, "autofill", false, "automatically fill shifts that are not scheduled yet")
	fs.StringVar(&startDate, "start", startDate, "starting with date, with --schedule")
	fs.IntVar(&numShifts, "shifts", numShifts, "number of shifts to scan, with --schedule")
	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "rotation forecast <rotation-name>")
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
