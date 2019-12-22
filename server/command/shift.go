// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) shift(parameters ...string) (string, error) {
	subcommands := map[string]func(...string) (string, error){
		"commit":       c.commitShift,
		"fill":         c.fillShift,
		"finish":       c.finishShift,
		"list":         c.listShifts,
		"open":         c.openShift,
		"start":        c.startShift,
		"debug-delete": c.deleteShift,
	}
	errUsage := errors.Errorf("Invalid subcommand. Usage:\n"+
		"- `%s shift open|fill|commit|start|finish [<rotation-name>] [flags...]` \n"+
		"\n"+
		"Use `%s rotation subcommand --help` for more information.\n",
		config.CommandTrigger, config.CommandTrigger)

	if len(parameters) == 0 {
		return "", errUsage
	}

	f := subcommands[parameters[0]]
	if f == nil {
		return "", errUsage
	}

	return f(parameters[1:]...)
}

func (c *Command) commitShift(parameters ...string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.CommitShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Committed shift #%v", shiftNumber), nil
		})
}

func (c *Command) fillShift(parameters ...string) (string, error) {
	return "", nil
}

func (c *Command) finishShift(parameters ...string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.FinishShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished shift #%v", shiftNumber), nil
		})
}

func (c *Command) listShifts(parameters ...string) (string, error) {
	start := time.Now().Format(api.DateFormat)
	numShifts := 3
	fs := flag.NewFlagSet("listShift", flag.ContinueOnError)
	fs.StringVarP(&start, "start", "s", start, fmt.Sprintf("Starting at, formatted as %s.", start))
	fs.IntVarP(&numShifts, "number", "n", 3, "Number of shifts to list")

	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "shift open <rotation-name>")
	if err != nil {
		return "", err
	}

	shifts, err := c.API.ListShifts(rotation, start, numShifts)
	if err != nil {
		return "", err
	}
	return utils.JSONBlock(shifts), nil
}

func (c *Command) openShift(parameters ...string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, err := c.API.OpenShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shift), nil
		})
}

func (c *Command) startShift(parameters ...string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.StartShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished shift #%v", shiftNumber), nil
		})
}

func (c *Command) deleteShift(parameters ...string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.DebugDeleteShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Deleted shift #%v", shiftNumber), nil
		})
}

func (c *Command) doShift(parameters []string, initf func(*pflag.FlagSet),
	updatef func(*api.Rotation, int) (string, error)) (string, error) {

	var start string
	var shiftNumber int
	fs := flag.NewFlagSet("doShift", flag.ContinueOnError)
	fs.StringVarP(&start, "start", "s", "", fmt.Sprintf("A date that would be in the shift, formatted as %s.", api.DateFormat))
	fs.IntVarP(&shiftNumber, "number", "n", -1, "Shift number")
	if initf != nil {
		initf(fs)
	}
	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "shift open <rotation-name>")
	if err != nil {
		return "", err
	}

	shiftNumber, err = c.shiftNumber(rotation, shiftNumber, start)
	if err != nil {
		return "", err
	}

	return updatef(rotation, shiftNumber)
}

func (c *Command) shiftNumber(rotation *api.Rotation, shiftNumber int, start string) (int, error) {
	if shiftNumber != -1 {
		if start != "" {
			return 0, errors.New("**Must specify --start or --number**")
		}
	} else {
		startTime, err := time.Parse(api.DateFormat, start)
		if err != nil {
			return 0, err
		}
		shiftNumber, err = rotation.ShiftNumberForTime(startTime)
		if err != nil {
			return 0, err
		}
	}
	return shiftNumber, nil
}
