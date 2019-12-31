// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func (c *Command) transitionShift(parameters []string) (string, error) {
	var commit, start, finish bool
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&commit, flagCommit, false, "Commit shift")
			fs.BoolVar(&start, flagStart, false, "Start shift")
			fs.BoolVar(&finish, flagFinish, false, "Finish shift")
		},
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			var err error
			switch {
			case commit && !start && !finish:
				err = c.API.CommitShift(rotation, shiftNumber)
			case start && !commit && !finish:
				err = c.API.StartShift(rotation, shiftNumber)
			case finish && !start && !commit:
				err = c.API.FinishShift(rotation, shiftNumber)
			case !finish && !start && !commit:
				return c.flagUsage(fs), errors.New("one of `--commit`, `--start`, `--finish` must be specified")
			default:
				return c.flagUsage(fs), errors.New("only one of `--commit`, `--start`, `--finish` can be specified")
			}
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Transitioned shift #%v", shiftNumber), nil
		})
}

func (c *Command) debugDeleteShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.DebugDeleteShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Deleted shift #%v", shiftNumber), nil
		})
}
