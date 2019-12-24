// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) add(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandRotation: c.addRotation,
		commandShift:    c.addShift,
		commandSkill:    c.addSkill,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `add rotation|shift|skill`")
}
