// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) userUnavailable(parameters []string) (string, error) {
	actingUser, err := c.SL.ActingUser()
	if err != nil {
		return "", err
	}

	fs := newFS()
	jsonOut := fJSON(fs)
	clear := fClear(fs)
	start, finish := fStartFinish(fs, actingUser)
	if err != nil {
		return "", err
	}
	err = fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	users, err := c.loadUsernames(fs.Args())
	if err != nil {
		return "", err
	}
	interval := types.Interval{
		Start:  *start,
		Finish: *finish,
	}

	if *clear {
		err = c.SL.ClearCalendar(users, interval)
		if err != nil {
			return "", err
		}
		if *jsonOut {
			return md.JSONBlock(users), nil
		}
		return fmt.Sprintf("cleared %s to %s from %s", start, finish, users.Markdown()), nil
	}

	u := sl.NewUnavailable(sl.ReasonPersonal, interval)
	err = c.SL.AddToCalendar(users, u)
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(users), nil
	}
	return fmt.Sprintf("Added %s to %s", actingUser.MarkdownUnavailable(u), users.Markdown()), nil
}
