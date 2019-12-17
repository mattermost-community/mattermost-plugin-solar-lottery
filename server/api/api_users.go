// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Users interface {
	// Returns the acting user, fully expanded
	GetActingUser() (*User, error)

	// Operations on lists of Mattermost usernames
	// Loads a list of Mattermost usernames as fully-expanded User's
	AddSkillToUsers(mattermostUsernames, skillName string, level int) error
	DeleteSkillFromUsers(mattermostUsernames, skillName string) error
	LoadMattermostUsers(mattermostUsernames string) (UserMap, error)

	// LoadStoredUser(s) loads only the store.User, leaves MattermostUser as nil
	// LoadStoredUser(mattermostUserID string) (*User, error)
	LoadStoredUsers(mattermostUserIDs store.IDMap) (UserMap, error)
}

func (api *api) GetActingUser() (*User, error) {
	err := api.Filter(withActingUser)
	if err != nil {
		return nil, err
	}
	return api.actingUser, nil
}

func (api *api) AddSkillToUsers(mattermostUsernames, skillName string, level int) error {
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
		MarkdownUser(api.actingUser), MarkdownSkillLevel(skillName, level), MarkdownUserMapWithSkills(api.users))
	return nil
}

func (api *api) DeleteSkillFromUsers(mattermostUsernames, skillName string) error {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
		withValidSkillName(skillName),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.AddSkillToUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"Skill":               skillName,
	})

	for _, user := range api.users {
		_, err = api.updateUserSkill(user, skillName, 0)
		if err != nil {
			return err
		}
	}

	logger.Infof("%s removed skill %s from %s.",
		MarkdownUser(api.actingUser), skillName, MarkdownUserMapWithSkills(api.users))
	return nil
}

func (api *api) LoadMattermostUsers(mattermostUsernames string) (UserMap, error) {
	users := UserMap{}
	names := strings.Split(mattermostUsernames, ",")
	for _, name := range names {
		if strings.HasPrefix(name, "@") {
			name = name[1:]
		}
		mmuser, err := api.PluginAPI.GetMattermostUserByUsername(name)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load %s", name)
		}
		user, _, err := api.loadOrMakeStoredUser(mmuser.Id)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load User %s", name)
		}
		user.MattermostUser = mmuser

		users[mmuser.Id] = user
	}
	return users, nil
}

func (api *api) LoadStoredUsers(ids store.IDMap) (UserMap, error) {
	users := UserMap{}
	for id := range ids {
		u, err := api.UserStore.LoadUser(id)
		if err != nil {
			return nil, err
		}
		users[u.MattermostUserID] = &User{
			User: u,
		}
	}
	return users, nil
}
