// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) rotation(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		// commandAutopilot:   c.autopilotRotation,
		// commandForecast:    c.forecastRotation,
		// commandGuess:       c.guessRotation,
		commandNew:         c.rotationNew,
		commandArchive:     c.rotationArchive,
		commandDebugDelete: c.rotationDebugDelete,
		commandParam:       c.rotationParam,
		commandList:        c.rotationList,
		commandShow:        c.rotationShow,
	}

	return c.handleCommand(subcommands, parameters)
}
