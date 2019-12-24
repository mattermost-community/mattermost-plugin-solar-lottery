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
	AddSkillToUsers(mattermostUsernames, skillName string, level Level) error
	DeleteSkillFromUsers(mattermostUsernames, skillName string) error
	VolunteerUsers(mattermostUsernames string, rotation *Rotation, shiftNumber int) error

	// LoadStoredUsers loads only the store.User, leaves MattermostUser as nil
	LoadStoredUsers(mattermostUserIDs store.IDMap) (UserMap, error)
	LoadMattermostUsers(mattermostUsernames string) (UserMap, error)
}

var ErrUserAlreadyInShift = errors.New("user is already in shift")

func (api *api) GetActingUser() (*User, error) {
	err := api.Filter(withActingUser)
	if err != nil {
		return nil, err
	}
	return api.actingUser, nil
}

func (api *api) AddSkillToUsers(mattermostUsernames, skillName string, level Level) error {
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

func (api *api) VolunteerUsers(mattermostUsernames string, rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.VolunteerUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"RotationID":          rotation.RotationID,
		"ShiftNumber":         shiftNumber,
	})

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return errors.Errorf("failed to load shift %v for rotation %s", shiftNumber, rotation.RotationID)
	}
	if shift.ShiftStatus != store.ShiftStatusOpen {
		return errors.Errorf("can't volunteer for a status which is %s, must be Open", shift.ShiftStatus)
	}

	volunteered := UserMap{}
	for _, user := range api.users {
		// TODO error if the shift has no vacancies? Or allow volunteering above the limit, and let the committer choose?
		if shift.Shift.MattermostUserIDs[user.MattermostUserID] != "" {
			return ErrUserAlreadyInShift
		}
		shift.Shift.MattermostUserIDs[user.MattermostUserID] = store.NotEmpty
		shift.Users[user.MattermostUserID] = user
		volunteered[user.MattermostUserID] = user
	}

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return err
	}

	api.messageShiftVolunteers(volunteered, rotation, shiftNumber, shift)
	logger.Infof("%s volunteered %s to %s shift %s.",
		MarkdownUser(api.actingUser), MarkdownUserMapWithSkills(volunteered), MarkdownRotation(rotation), MarkdownShift(shiftNumber, shift))
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
