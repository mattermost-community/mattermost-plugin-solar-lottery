// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const eUser = "user_"

type Users interface {
	LoadMattermostUsername(username string) (*User, error)
	Disqualify(users UserMap, skillName string) error
	Qualify(users UserMap, skillLevel SkillLevel) error
}

func (sl *sl) ActingUser() (*User, error) {
	err := sl.Filter(withActingUser)
	if err != nil {
		return nil, err
	}
	return sl.actingUser, nil
}

func (sl *sl) loadStoredUsers(ids *types.Set) (UserMap, error) {
	users := UserMap{}
	err := ids.ForEachWithError(func(id string) error {
		var user User
		err := sl.Store.Entity(KeyUser).Load(id, &user)
		if err != nil {
			return err
		}
		users[id] = &user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (sl *sl) LoadMattermostUsername(username string) (*User, error) {
	err := sl.Filter(withActingUserExpanded)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(username, "@") {
		username = username[1:]
	}
	mmuser, err := sl.PluginAPI.GetMattermostUserByUsername(username)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to load %s", username)
	}
	user, _, err := sl.loadOrMakeUser(mmuser.Id)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to load User %s", username)
	}
	user.mattermostUser = mmuser
	return user, nil
}

func (sl *sl) Qualify(users UserMap, skillLevel SkillLevel) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "Qualify",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"Users":          users.String(),
		"SkillLevel":     skillLevel,
	})

	err = sl.AddKnownSkill(skillLevel.Skill)
	if err != nil {
		return err
	}
	for _, user := range users {
		err = sl.updateUserSkill(user, skillLevel)
		if err != nil {
			return err
		}
	}

	logger.Infof("%s added skill %s to %s.",
		sl.actingUser.Markdown(), skillLevel, users.Markdown())
	return nil
}

func (sl *sl) Disqualify(users UserMap, skillName string) error {
	err := sl.Filter(
		withActingUserExpanded,
		withValidSkillName(skillName),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":   "sl.AddSkillToUsers",
		"ActingUser": sl.actingUser.MattermostUsername(),
		"Users":      users.String(),
		"Skill":      skillName,
	})

	for _, user := range users {
		err = sl.updateUserSkill(user, SkillLevel{
			Skill: skillName,
			Level: Level(0),
		})
		if err != nil {
			return err
		}
	}

	logger.Infof("%s removed skill %s from %s.",
		sl.actingUser.Markdown(), skillName, users.Markdown())
	return nil
}

func (sl *sl) ExpandUser(user *User) error {
	if user.mattermostUser != nil {
		return nil
	}
	mattermostUser, err := sl.PluginAPI.GetMattermostUser(user.MattermostUserID)
	if err != nil {
		return err
	}
	user.mattermostUser = mattermostUser

	loc, err := time.LoadLocation(mattermostUser.GetPreferredTimezone())
	if err != nil {
		return err
	}
	user.location = loc
	return nil
}

func (sl *sl) loadOrMakeUser(mattermostUserID string) (*User, bool, error) {
	var user User
	err := sl.Store.Entity(KeyUser).Load(mattermostUserID, &user)
	if err == kvstore.ErrNotFound {
		newUser, newErr := sl.storeUserWelcomeNew(NewUser(mattermostUserID))
		return newUser, true, newErr
	}
	if err != nil {
		return nil, false, err
	}
	return &user, false, nil
}

// storeUserWelcomeNew checks if the user being stored is new, and welcomes the user.
// note that it can be used inside of filters, so it must not use filters itself,
//  nor assume that any runtime values have been filled.
func (sl *sl) storeUserWelcomeNew(orig *User) (*User, error) {
	user, err := sl.storeUser(orig)
	if err != nil {
		return nil, err
	}
	if orig.PluginVersion == "" {
		sl.messageWelcomeNewUser(user)
	}
	return user, nil
}

func (sl *sl) storeUser(orig *User) (*User, error) {
	user := orig.Clone()
	user.PluginVersion = sl.Config().PluginVersion
	err := sl.Store.Entity(KeyUser).Store(user.MattermostUserID, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (sl *sl) updateUserSkill(user *User, skillLevel SkillLevel) error {
	s, l := skillLevel.Skill, skillLevel.Level
	if user.SkillLevels[s] == int64(l) {
		// nothing to do
		sl.Logger.Debugf("nothing to do for user %s, already is %s", user.Markdown(), skillLevel)
		return nil
	}

	if l == 0 {
		_, ok := user.SkillLevels[s]
		if ok {
			delete(user.SkillLevels, s)
		}
	} else {
		user.SkillLevels[s] = int64(l)
	}
	user, err := sl.storeUserWelcomeNew(user)
	if err != nil {
		return err
	}
	sl.Logger.Debugf("%s updated to %s", user.Markdown(), skillLevel)
	return nil
}
