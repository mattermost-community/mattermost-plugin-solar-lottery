// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func (c *Command) rotation(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandAutopilot:   c.autopilotRotation,
		commandAdd:         c.addRotation,
		commandArchive:     c.archiveRotation,
		commandDebugDelete: c.debugDeleteRotation,
		commandForecast:    c.forecastRotation,
		commandGuess:       c.guessRotation,
		commandJoin:        c.joinRotation,
		commandLeave:       c.leaveRotation,
		commandList:        c.listRotations,
		commandNeed:        c.rotationNeed,
		commandShow:        c.showRotation,
		commandUpdate:      c.updateRotation,
	}

	return c.handleCommand(subcommands, parameters)
}

func withRotationFlags(fs *pflag.FlagSet, rotationID, rotationName *string) {
	fs.StringVar(rotationID, flagRotationID, "", "specify rotation ID")
	fs.StringVarP(rotationName, flagRotation, flagPRotation, "", "specify rotation name")
}

func (c *Command) parseRotationFlags(id, name string) (rotationID string, err error) {
	switch {
	case id == "" && name == "":
		return "", errors.New("rotation is not specified")

	case id != "" && name != "":
		return "", errors.New("rotation is specified multiple times")

	case id != "":
		return id, nil

	}
	//  name != "":
	rotationIDs, err := c.API.ResolveRotationName(name)
	if err != nil {
		return "", err
	}
	if len(rotationIDs) != 1 {
		return "", errors.Errorf("name %s is ambigous, please use --%s with one of %s", name, flagRotation, rotationIDs)
	}
	return rotationIDs[0], nil
}

func newRotationFlagSet(rotationID, rotationName *string) *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withRotationFlags(fs, rotationID, rotationName)
	return fs
}
