// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) show(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandRotation: c.showRotation,
		commandShow:     c.showShift,
		commandUser:     c.showUser,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `show rotation|shift|user`.")
}
