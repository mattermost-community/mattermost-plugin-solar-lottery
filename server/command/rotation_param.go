// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationParam(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandGrace:  c.rotationParamGrace,
		commandMax:    c.rotationParamMax,
		commandMin:    c.rotationParamMin,
		commandShift:  c.rotationParamShift,
		commandTicket: c.rotationParamTicket,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) rotationParamMin(parameters []string) (md.MD, error) {
	return c.rotationParamMinMax(true, parameters)
}

func (c *Command) rotationParamMax(parameters []string) (md.MD, error) {
	return c.rotationParamMinMax(false, parameters)
}

func (c *Command) rotationParamMinMax(min bool, parameters []string) (md.MD, error) {
	c.withFlagRotation()
	skillLevel := c.withFlagSkillLevel()
	count := c.withFlagCount()
	clear := c.withFlagClear()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			needsToUpdate := r.TaskMaker.Limit
			if min {
				needsToUpdate = r.TaskMaker.Require
			}
			if *clear {
				needsToUpdate.Delete(skillLevel.AsID())
			} else {
				needsToUpdate.SetCountForSkillLevel(*skillLevel, int64(*count))
			}
			return nil
		}))
}

func (c *Command) rotationParamGrace(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	dur := c.withFlagDuration()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.TaskMaker.Grace = *dur
			return nil
		}))
}

func (c *Command) rotationParamShift(parameters []string) (md.MD, error) {
	actingUser, err := c.SL.ActingUser()
	if err != nil {
		return "", err
	}

	c.withFlagRotation()
	period := c.withFlagPeriod()
	start := c.withFlagStart(actingUser)
	err = c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.TaskMaker.Type = sl.ShiftMaker
			r.TaskMaker.ShiftPeriod = *period
			r.TaskMaker.ShiftStart = *start
			return nil
		}))
}

func (c *Command) rotationParamTicket(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.TaskMaker.Type = sl.TicketMaker
			return nil
		}))
}
