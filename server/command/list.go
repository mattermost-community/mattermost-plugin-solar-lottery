// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) list(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandRotation: c.listRotations,
		commandShift:    c.listShifts,
		commandSkill:    c.listSkills,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `list rotation|shift|skill`")
}
