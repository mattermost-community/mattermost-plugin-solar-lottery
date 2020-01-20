// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) shift(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandOpen:        c.openShift,
		commandDebugDelete: c.debugDeleteShift,
		commandFill:        c.fillShift,
		commandJoin:        c.joinShift,
		commandList:        c.listShifts,
		commandStart:       c.startShift,
		commandFinish:      c.finishShift,
		// commandLeave:       c.leaveShift,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) doShift(parameters []string,
	initF func(fs *pflag.FlagSet),
	doF func(*pflag.FlagSet, *api.Rotation, int) (string, error)) (string, error) {
	var rotationID, rotationName string
	var shiftNumber int
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.IntVarP(&shiftNumber, flagShift, flagPShift, 0, "(starting) shift number")
	if initF != nil {
		initF(fs)
	}
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

	return doF(fs, rotation, shiftNumber)
}

func (c *Command) openShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, err := c.API.OpenShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shift), nil
		})
}

func (c *Command) startShift(parameters []string) (string, error) {
	return c.doShift(parameters,
		nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			_, err := c.API.StartShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Started %s", rotation.ShiftRef(shiftNumber)), nil
		})
}

func (c *Command) finishShift(parameters []string) (string, error) {
	return c.doShift(parameters,
		nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			_, err := c.API.FinishShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished %s", rotation.ShiftRef(shiftNumber)), nil
		})
}

func (c *Command) debugDeleteShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.DebugDeleteShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Deleted shift #%v", shiftNumber), nil
		})
}
