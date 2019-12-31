// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) user(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandDisqualify:  c.disqualifyUsers,
		commandQualify:     c.qualifyUsers,
		commandShow:        c.showUser,
		commandUnavailable: c.userUnavailable,
	}
	return c.handleCommand(subcommands, parameters)
}
