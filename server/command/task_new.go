// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) newTask(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandTicket: c.newTicketTask,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) newTicketTask(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	summary := fSummary(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	t, err := c.SL.MakeTicket(rotationID, *summary, "")
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(t), nil
	}
	return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), rotationID), nil
}
