// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) userUnavailable(parameters []string) (md.MD, error) {
	clear := c.withFlagClear()
	start, err := c.withFlagStart()
	if err != nil {
		return "", err
	}
	finish, err := c.withFlagFinish()
	if err != nil {
		return "", err
	}
	err = c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	mattermostUserIDs, err := c.resolveUsernames(c.fs.Args())
	if err != nil {
		return "", err
	}
	interval := types.NewInterval(*start, *finish)

	if *clear {
		return c.normalOut(
			c.SL.ClearCalendar(sl.InClearCalendar{
				MattermostUserIDs: mattermostUserIDs,
				Interval:          interval,
			}))
	}

	return c.normalOut(
		c.SL.AddToCalendar(sl.InAddToCalendar{
			MattermostUserIDs: mattermostUserIDs,
			Unavailable:       sl.NewUnavailable(sl.ReasonPersonal, interval),
		}))
}
