// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) rotation(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandArchive:     c.rotationArchive,
		commandAutopilot:   c.rotationAutopilot,
		commandDebugDelete: c.rotationDebugDelete,
		commandList:        c.rotationList,
		commandNew:         c.rotationNew,
		commandParam:       c.rotationParam,
		commandShow:        c.rotationShow,
	}

	return c.handleCommand(subcommands, parameters)
}
