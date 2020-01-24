// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

func (sl *solarLottery) GetActingUser() (*User, error) {
	err := sl.Filter(withActingUser)
	if err != nil {
		return nil, err
	}
	return sl.actingUser, nil
}

func (sl *solarLottery) LoadStoredUsers(ids store.IDMap) (UserMap, error) {
	users := UserMap{}
	for id := range ids {
		u, err := sl.UserStore.LoadUser(id)
		if err != nil {
			return nil, err
		}
		users[u.MattermostUserID] = &User{
			User: u,
		}
	}
	return users, nil
}

func (sl *solarLottery) LoadMattermostUsers(mattermostUsernames string) (UserMap, error) {
	err := sl.Filter(withActingUserExpanded)
	if err != nil {
		return nil, err
	}

	if mattermostUsernames == "" {
		return UserMap{
			sl.actingMattermostUserID: sl.actingUser,
		}, nil
	}

	users := UserMap{}
	names := strings.Split(mattermostUsernames, ",")
	for _, name := range names {
		if strings.HasPrefix(name, "@") {
			name = name[1:]
		}
		mmuser, err := sl.PluginAPI.GetMattermostUserByUsername(name)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load %s", name)
		}
		user, _, err := sl.loadOrMakeStoredUser(mmuser.Id)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load User %s", name)
		}
		user.MattermostUser = mmuser
		users[mmuser.Id] = user
	}
	return users, nil
}

func (sl *solarLottery) Qualify(mattermostUsernames, skillName string, level Level) error {
	err := sl.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":            "sl.AddSkillToUsers",
		"ActingUsername":      sl.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"Skill":               skillName,
		"Level":               level,
	})

	err = sl.AddSkill(skillName)
	if err != nil && err != ErrAlreadyExists {
		return err
	}

	for _, user := range sl.users {
		err = sl.updateUserSkill(user, skillName, level)
		if err != nil {
			return err
		}
	}

	logger.Infof("%s added skill %s to %s.",
		sl.actingUser.Markdown(), MarkdownSkillLevel(skillName, level), sl.users.MarkdownWithSkills())
	return nil
}

func (sl *solarLottery) Disqualify(mattermostUsernames, skillName string) error {
	err := sl.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
		withValidSkillName(skillName),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":            "sl.AddSkillToUsers",
		"ActingUsername":      sl.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"Skill":               skillName,
	})

	for _, user := range sl.users {
		err = sl.updateUserSkill(user, skillName, 0)
		if err != nil {
			return err
		}
	}

	logger.Infof("%s removed skill %s from %s.",
		sl.actingUser.Markdown(), skillName, sl.users.MarkdownWithSkills())
	return nil
}
