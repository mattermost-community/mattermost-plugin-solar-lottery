// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) shift(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandOpen:        c.openShift,
		commandDebugDelete: c.debugDeleteShift,
		commandFill:        c.fillShift,
		commandJoin:        c.joinShift,
		commandList:        c.listShifts,
		commandTransition:  c.transitionShift,
		// commandLeave:       c.leaveShift,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) doShift(parameters []string,
	initF func(fs *pflag.FlagSet),
	doF func(*pflag.FlagSet, *api.Rotation, int) (string, error)) (string, error) {
	var rotationID, rotationName string
	var shiftNumber int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.IntVarP(&shiftNumber, flagShift, flagPShift, 0, "(starting) shift number")
	if initF != nil {
		initF(fs)
	}
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

	return doF(fs, rotation, shiftNumber)
}
