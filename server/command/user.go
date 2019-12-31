// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) user(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandDisqualify:  c.disqualifyUsers,
		commandQualify:     c.qualifyUsers,
		commandShow:        c.showUser,
		commandUnavailable: c.userUnavailable,
	}
	return c.handleCommand(subcommands, parameters,
		"Usage: `user disqualify|qualify|show|unavailable`. Use `user subcommand --help` for more information.")
}

func (c *Command) showUser(parameters []string) (string, error) {
	usernames := ""
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
	}

	users, err := c.API.LoadMattermostUsers(usernames)
	if err != nil {
		return "", err
	}
	return utils.JSONBlock(users), nil
}

func (c *Command) userUnavailable(parameters []string) (string, error) {
	var typ, usernames, start, end string
	var clear bool
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
	fs.StringVarP(&start, flagStart, flagPStart, "", "start of the unavailability")
	fs.StringVarP(&end, flagEnd, flagPEnd, "", "end of unavailability (last day)")
	fs.BoolVar(&clear, flagClear, false, "clear all overlapping events")
	fs.StringVar(&typ, flagType, store.EventTypeOther, "event type")
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
	}

	endTime, err := time.Parse(api.DateFormat, end)
	if err != nil {
		return "", err
	}
	endTime = endTime.Add(time.Hour * 24) // start of next day
	end = endTime.Format(api.DateFormat)

	if clear {
		err = c.API.DeleteEventsFromUsers(usernames, start, end)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("cleared %s to %s from %s", start, end, usernames), nil
	} else {
		event := store.Event{
			Start: start,
			End:   endTime.Format(api.DateFormat),
			Type:  typ,
		}
		err = c.API.AddEventToUsers(usernames, event)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Added %s to %s", api.MarkdownEvent(event), usernames), nil
	}
}
