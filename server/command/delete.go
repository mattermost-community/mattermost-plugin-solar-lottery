// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) delete(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandDebugRotation: c.debugDeleteRotation,
		commandDebugShift:    c.debugDeleteShift,
		commandSkill:         c.deleteSkill,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `delete debug-rotation|debug-shift|skill`")
}
