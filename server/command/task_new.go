// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) taskNew(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandTicket: c.newTicketTask,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) newTicketTask(parameters []string) (md.MD, error) {
	fs := c.assureFS()
	c.withFlagRotation()
	summary := c.withFlagSummary()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(c.SL.MakeTicket(sl.InMakeTicket{
		RotationID: rotationID,
		Summary:    *summary,
	}))
}
