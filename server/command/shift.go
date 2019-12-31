// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) shift(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandAdd:         c.addShift,
		commandCommit:      c.commitShift,
		commandDebugDelete: c.debugDeleteShift,
		commandFill:        c.fillShift,
		commandFinish:      c.finishShift,
		commandJoin:        c.joinShift,
		// commandLeave:       c.leaveShift,
		commandStart: c.startShift,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `shift add|commit|fill|finish|start]`. Use `rotation subcommand --help` for more information.")
}

func (c *Command) commitShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, start int) (string, error) {
			err := c.API.CommitShift(rotation, start)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Committed shift #%v", start), nil
		})
}

func (c *Command) fillShift(parameters []string) (string, error) {
	return "TODO", nil
}

func (c *Command) finishShift(parameters []string) (string, error) {
	return c.doShift(parameters,
		nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.FinishShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished shift #%v", shiftNumber), nil
		})
}

func (c *Command) listShifts(parameters []string) (string, error) {
	numShifts := 0
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.IntVarP(&numShifts, flagNumber, flagPNumber, 3, "Number of shifts to list")
		},
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			shifts, err := c.API.ListShifts(rotation, shiftNumber, numShifts)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shifts), nil
		})
}

func (c *Command) addShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, err := c.API.OpenShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shift), nil
		})
}

func (c *Command) showShift(parameters []string) (string, error) {
	return c.doShift(
		parameters,
		nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, err := c.API.OpenShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shift), nil
		})
}

func (c *Command) startShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.StartShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished shift #%v", shiftNumber), nil
		})
}

func (c *Command) debugDeleteShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.DebugDeleteShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Deleted shift #%v", shiftNumber), nil
		})
}

func (c *Command) doShift(parameters []string,
	initF func(fs *pflag.FlagSet),
	doF func(*api.Rotation, int) (string, error)) (string, error) {
	var rotationID, rotationName string
	var start int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.IntVarP(&start, flagStart, flagPStart, 0, "starting shift number")
	// fs.IntVarP(&numShifts, flagNumber, flagPNumber, -1, "number of shifts")
	if initF != nil {
		initF(fs)
	}
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

	return doF(rotation, start)
}
