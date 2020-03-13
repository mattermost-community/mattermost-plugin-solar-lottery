// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) task(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandParam:  c.taskParam,
		commandNew:    c.newTask,
		commandShow:   c.showTask,
		commandAssign: c.assignTask,
		commandFill:   c.fillTask,
	}

	return c.handleCommand(subcommands, parameters)
}
