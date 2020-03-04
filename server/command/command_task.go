// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) task(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandParam: c.taskParam,
		commandNew:   c.newTask,
	}

	return c.handleCommand(subcommands, parameters)
}
