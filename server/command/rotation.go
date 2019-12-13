// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) rotation(parameters ...string) (string, error) {
	subcommands := map[string]func(...string) (string, error){
		"list":   c.listRotations,
		"show":   c.showRotation,
		"delete": c.deleteRotation,
		"add":    c.addRotation,
		"update": c.updateRotation,
	}
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}

	f := subcommands[parameters[0]]
	if f == nil {
		return "", errors.New("invalid syntax TODO")
	}

	return f(parameters[1:]...)
}

func (c *Command) showRotation(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]

	rr, err := c.API.ListRotations()
	if err != nil {
		return "", err
	}
	r := rr[rotationName]
	if r == nil {
		return "", store.ErrNotFound
	}

	s := flag.NewFlagSet("rotation", flag.ContinueOnError)
	schedule := false
	autofill := false
	start, err := api.ShiftNumber(r, time.Now())
	if err != nil {
		return "", err
	}
	numShifts := 12
	s.BoolVar(&schedule, "schedule", false, "display future schedule")
	s.BoolVar(&autofill, "autofill", false, "automatically fill shifts that are not scheduled yet")
	s.IntVar(&start, "start", start, "starting shift to display, with --schedule")
	s.IntVar(&numShifts, "shifts", numShifts, "number of shifts forward, with --schedule")
	err = s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}

	if !schedule {
		return utils.JSONBlock(r), nil
	}

	shifts, err := c.API.CrystalBall(rotationName, start, numShifts, autofill)
	if err != nil {
		return "", err
	}

	return utils.JSONBlock(shifts), nil
}

func (c *Command) listRotations(parameters ...string) (string, error) {
	rr, err := c.API.ListRotations()
	if err != nil {
		return "", err
	}
	out := ""
	for name := range rr {
		out += fmt.Sprintf("- %s\n", name)
	}
	return out, nil
}

func (c *Command) deleteRotation(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]
	err := c.API.DeleteRotation(rotationName)
	if err != nil {
		return "", err
	}
	return "Deleted rotation " + rotationName, nil
}

func (c *Command) addRotation(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]

	r := store.NewRotation(rotationName)
	s := flag.NewFlagSet("rotation", flag.ContinueOnError)
	s.StringVar(&r.Period, "period", "m", "rotation period 'w', '2w', or 'm'")
	s.StringVar(&r.Start, "start", "", "rotation starts on")
	s.IntVar(&r.MinBetweenShifts, "min-between", 1, "minimum number of shifts between being elected")
	s.IntVar(&r.Size, "size", 1, "number of people in the shift")
	err := s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}

	//TODO input validation

	r, err = c.API.AddRotation(r)
	if err != nil {
		return "", err
	}

	return "Added rotation " + utils.JSONBlock(r), nil
}

func (c *Command) updateRotation(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]

	s := flag.NewFlagSet("rotation", flag.ContinueOnError)
	var noNeedName, needName, needSkill, needLevel, period, start string
	var minBetween, size, needCount int
	s.StringVar(&period, "period", "", "rotation period 'w', '2w', or 'm'")
	s.StringVar(&start, "start", "", "rotation starts on")
	s.StringVar(&noNeedName, "no-need", "", "remove a need from the rotation")
	s.StringVar(&needName, "need", "", "update rotation's needs")
	s.StringVar(&needSkill, "skill", "", "if used with --need, indicates the needed skill")
	s.StringVar(&needLevel, "level", "", "if used with --need, indicates the needed skill level")
	s.IntVar(&needCount, "count", 0, "if used with --need, indicates the needed headcount")
	s.IntVar(&minBetween, "min-between", 0, "minimum number of periods between shifts")
	s.IntVar(&size, "size", 0, "number of people in the shift")
	err := s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}

	if needName != "" && len(needLevel)+len(needSkill)+needCount == 0 {
		return "", errors.New("--need requires skill, level, and count to be specified")
	}
	if needName != "" && noNeedName != "" {
		return "", errors.New("--need and --no-need can not be used in the same command")
	}

	// TODO more input validation

	r, err := c.API.UpdateRotation(rotationName, func(r *store.Rotation) error {
		if period != "" {
			r.Period = period
		}
		if start != "" {
			r.Start = start
		}
		if minBetween != 0 {
			r.MinBetweenShifts = minBetween
		}
		if size != 0 {
			r.Size = size
		}
		if needName != "" {
			level := 0
			level, err = api.ParseLevel(needLevel)
			if err != nil {
				return err
			}
			c.API.ChangeRotationNeed(r, needName, needSkill, int(level), needCount)
		}
		if noNeedName != "" {
			err = c.API.RemoveRotationNeed(r, noNeedName)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return "Updated rotation " + utils.JSONBlock(r), nil
}
