// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) DeleteEvents(mattermostUsernames string, startDate, endDate string) error {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.AddSkillToUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"StartDate":           startDate,
		"EndDate":             endDate,
	})

	for _, user := range api.users {
		intervalStart, intervalEnd, err := ParseDatePair(startDate, endDate)
		if err != nil {
			return err
		}

		_, err = user.overlapEvents(intervalStart, intervalEnd, true)
		if err != nil {
			return errors.WithMessagef(err, "failed to remove events from %s to %s", startDate, endDate)
		}

		_, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return errors.WithMessagef(err, "failed to update user %s", api.MarkdownUser(user))
		}
	}

	logger.Infof("%s deleted events from %s to %s from users %s.",
		api.MarkdownUser(api.actingUser), startDate, endDate, api.MarkdownUsersWithSkills(api.users))
	return nil
}
