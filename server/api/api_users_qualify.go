// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) Qualify(mattermostUsernames, skillName string, level Level) error {
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
		"Skill":               skillName,
		"Level":               level,
	})

	err = api.AddSkill(skillName)
	if err != nil && err != ErrSkillAlreadyExists {
		return err
	}

	for _, user := range api.users {
		_, err = api.updateUserSkill(user, skillName, level)
		if err != nil {
			return err
		}
	}

	logger.Infof("%s added skill %s to %s.",
		api.MarkdownUser(api.actingUser), api.MarkdownSkillLevel(skillName, level), api.MarkdownUsersWithSkills(api.users))
	return nil
}
