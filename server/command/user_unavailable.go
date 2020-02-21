// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) userUnavailable(parameters []string) (string, error) {
	var start, finish types.Time
	var clear bool
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.VarP(&start, flagStart, flagPStart, "start of the interval")
	fs.VarP(&finish, flagEnd, flagPEnd, "end of the interval")
	fs.BoolVar(&clear, flagClear, false, "mark as available by clearing all overlapping unavailability events")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	users, err := c.loadUsernames(fs.Args())
	if err != nil {
		return "", err
	}
	interval := types.Interval{
		Start:  start,
		Finish: finish,
	}

	if clear {
		err = c.SL.ClearCalendar(users, interval)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("cleared %s to %s from %s", start, finish, users.Markdown()), nil
	}

	u := sl.NewUnavailable(sl.ReasonPersonal, interval)
	err = c.SL.AddToCalendar(users, u)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Added %s to %s", u.Markdown(), users.Markdown()), nil
}
