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
		commandAdd:         c.addShift,
		commandCommit:      c.commitShift,
		commandDebugDelete: c.debugDeleteShift,
		commandFill:        c.fillShift,
		commandFinish:      c.finishShift,
		commandJoin:        c.joinShift,
		// commandLeave:       c.leaveShift,
		commandStart: c.startShift,
	}

	return c.handleCommand(subcommands, parameters)
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

	return doF(rotation, start)
}
