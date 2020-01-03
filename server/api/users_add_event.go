// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) AddEvent(mattermostUsernames string, event Event) error {
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
		"Event":               event,
	})

	err = api.addEventToUsers(api.users, event, true)
	if err != nil {
		return err
	}

	logger.Infof("%s added event %s to %s.",
		api.MarkdownUser(api.actingUser), MarkdownEvent(event), api.MarkdownUsersWithSkills(api.users))
	return nil
}

func (api *api) addEventToUsers(users UserMap, event Event, persist bool) error {
	for _, user := range users {
		user.AddEvent(event)

		if persist {
			_, err := api.storeUserWelcomeNew(user)
			if err != nil {
				return errors.WithMessagef(err, "failed to update user %s", api.MarkdownUser(user))
			}
		}
	}
	return nil
}
