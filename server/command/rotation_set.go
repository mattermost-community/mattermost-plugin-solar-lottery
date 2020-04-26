// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationSet(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandAutopilot: c.rotationSetAutopilot,
		commandFill:      c.rotationSetFill,
		commandLimit:     c.rotationSetLimit,
		commandRequire:   c.rotationSetRequire,
		commandTask:      c.rotationSetTask,
	}

	return c.handleCommand(subcommands, parameters)
}
