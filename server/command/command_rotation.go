// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) rotation(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		// commandAutopilot:   c.autopilotRotation,
		commandNew:         c.newRotation,
		commandArchive:     c.archiveRotation,
		commandDebugDelete: c.debugDeleteRotation,
		// commandForecast:    c.forecastRotation,
		// commandGuess:       c.guessRotation,
		commandJoin:  c.joinRotation,
		commandLeave: c.leaveRotation,
		commandList:  c.listRotations,
		// commandNeed:  c.rotationNeed,
		commandShow: c.showRotation,
	}

	return c.handleCommand(subcommands, parameters)
}
