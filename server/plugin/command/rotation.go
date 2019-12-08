// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func (c *Command) rotation(parameters ...string) (string, error) {
	subcommands := map[string]func(...string) (string, error){
		"list":   c.listRotations,
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

func (c *Command) listRotations(parameters ...string) (string, error) {
	rr, err := c.API.ListRotations()
	if err != nil {
		return "", err
	}
	return utils.JSONBlock(rr), nil
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
	period, start := "", ""
	minBetween, size := 0, 0
	s.StringVar(&period, "period", "", "rotation period 'w', '2w', or 'm'")
	s.StringVar(&start, "start", "", "rotation starts on")
	s.IntVar(&minBetween, "min-between", 0, "minimum number of periods between shifts")
	s.IntVar(&size, "size", 0, "number of people in the shift")
	err := s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}

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
		return nil
	})
	if err != nil {
		return "", err
	}

	return "Updated rotation " + utils.JSONBlock(r), nil
}
