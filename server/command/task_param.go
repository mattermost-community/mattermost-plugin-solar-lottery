// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) taskParam(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandMin:    c.taskParamMin,
		commandMax:    c.taskParamMax,
		commandGrace:  c.taskParamGrace,
		commandShift:  c.taskParamShift,
		commandTicket: c.taskParamTicket,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) taskParamMin(parameters []string) (md.MD, error) {
	return c.taskParamMinMax(true, parameters)
}

func (c *Command) taskParamMax(parameters []string) (md.MD, error) {
	return c.taskParamMinMax(false, parameters)
}

func (c *Command) taskParamMinMax(min bool, parameters []string) (md.MD, error) {
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
			needsToUpdate := r.TaskMaker.Max
			if min {
				needsToUpdate = r.TaskMaker.Min
			}
			if *clear {
				needsToUpdate.Delete(skillLevel.AsID())
			} else {
				needsToUpdate.SetCountForSkillLevel(*skillLevel, int64(*count))
			}
			return nil
		}))
}

func (c *Command) taskParamGrace(parameters []string) (md.MD, error) {
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

func (c *Command) taskParamShift(parameters []string) (md.MD, error) {
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

func (c *Command) taskParamTicket(parameters []string) (md.MD, error) {
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
