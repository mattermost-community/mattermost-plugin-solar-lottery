// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) user(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandDisqualify:  c.disqualifyUsers,
		commandQualify:     c.qualifyUsers,
		commandShow:        c.showUser,
		commandUnavailable: c.userUnavailable,
		// commandForecast:    c.userForecast,
		commandJoin:  c.joinRotation,
		commandLeave: c.leaveRotation,
	}

	return c.handleCommand(subcommands, parameters)
}
