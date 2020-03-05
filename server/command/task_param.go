// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) taskParam(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandMin:    c.taskParamMin,
		commandMax:    c.taskParamMax,
		commandGrace:  c.taskParamGrace,
		commandShift:  c.taskParamShift,
		commandTicket: c.taskParamTicket,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) taskParamMin(parameters []string) (string, error) {
	return c.taskParamMinMax(true, parameters)
}

func (c *Command) taskParamMax(parameters []string) (string, error) {
	return c.taskParamMinMax(false, parameters)
}

func (c *Command) taskParamMinMax(min bool, parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	skillLevel := fSkillLevel(fs)
	count := fCount(fs)
	clear := fClear(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	r, err := c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
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
	})
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), rotationID), nil
}

func (c *Command) taskParamGrace(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	dur := fDuration(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	r, err := c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
		r.TaskMaker.Grace = *dur
		return nil
	})
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), rotationID), nil
}

func (c *Command) taskParamShift(parameters []string) (string, error) {
	actingUser, err := c.SL.ActingUser()
	if err != nil {
		return "", err
	}

	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	period := fPeriod(fs)
	start := fStart(fs, actingUser)
	err = fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	r, err := c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
		r.TaskMaker.Type = sl.ShiftMaker
		r.TaskMaker.ShiftPeriod = *period
		r.TaskMaker.ShiftStart = *start
		return nil
	})
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), rotationID), nil
}

func (c *Command) taskParamTicket(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}

	r, err := c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
		r.TaskMaker.Type = sl.TicketMaker
		return nil
	})
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), rotationID), nil
}
