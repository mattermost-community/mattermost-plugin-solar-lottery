// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) rotation(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		// commandAutopilot:   c.autopilotRotation,
		// commandForecast:    c.forecastRotation,
		// commandGuess:       c.guessRotation,
		commandNew:         c.newRotation,
		commandArchive:     c.archiveRotation,
		commandDebugDelete: c.debugDeleteRotation,
		commandList:        c.listRotations,
		commandShow:        c.showRotation,
	}

	return c.handleCommand(subcommands, parameters)
}
