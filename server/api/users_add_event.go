// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) AddEvent(mattermostUsernames string, event store.Event) error {
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

	for _, user := range api.users {
		err = user.AddEvent(event)
		if err != nil {
			return errors.WithMessagef(err, "failed to add event %s to %s", MarkdownEvent(event), MarkdownUser(user))
		}

		_, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return errors.WithMessagef(err, "failed to update user %s", MarkdownUser(user))
		}
	}

	logger.Infof("%s added event %s to %s.",
		MarkdownUser(api.actingUser), MarkdownEvent(event), MarkdownUserMapWithSkills(api.users))
	return nil
}
