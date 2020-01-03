// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

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
		api.MarkdownUser(api.actingUser), MarkdownSkillLevel(skillName, level), api.MarkdownUsersWithSkills(api.users))
	return nil
}

func (api *api) updateUserSkill(user *User, skillName string, level Level) (*User, error) {
	if user.SkillLevels[skillName] == int(level) {
		// nothing to do
		api.Logger.Debugf("nothing to do for user %s, already has skill %s (%v)", api.MarkdownUser(user), skillName, level)
		return user, nil
	}

	if level == 0 {
		_, ok := user.SkillLevels[skillName]
		if !ok {
			return nil, errors.Errorf("%s does not have skill %s", api.MarkdownUser(user), skillName)
		}
		delete(user.SkillLevels, skillName)
	} else {
		user.SkillLevels[skillName] = int(level)
	}

	user, err := api.storeUserWelcomeNew(user)
	if err != nil {
		return nil, err
	}
	api.Logger.Debugf("%s (%v) skill updated user %s", skillName, level, api.MarkdownUser(user))
	return user, nil
}
