// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) taskNew(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandTicket: c.taskNewTicket,
		commandShift:  c.taskNewShift,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) taskNewTicket(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	summary := c.assureFS().String("summary", "", "task summary")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.CreateTicket(sl.InCreateTicket{
			RotationID: rotationID,
			Summary:    *summary,
			Time:       *c.now,
		}))
}

func (c *Command) taskNewShift(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	shiftNumber := c.assureFS().IntP("number", "n", 1, "shift number")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(c.SL.CreateShift(sl.InCreateShift{
		RotationID: rotationID,
		Number:     *shiftNumber,
		Time:       types.NewTime(time.Now()),
	}))
}
