// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) user(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandDisqualify:  c.userDisqualify,
		commandQualify:     c.userQualify,
		commandShow:        c.userShow,
		commandUnavailable: c.userUnavailable,
		commandJoin:        c.userJoin,
		commandLeave:       c.userLeave,
	}

	return c.handleCommand(subcommands, parameters)
}
