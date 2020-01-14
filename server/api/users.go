// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type User struct {
	*store.User

	// nil is assumed to be valid
	MattermostUser *model.User
}

func (user *User) Clone() *User {
	clone := *user
	clone.User = user.User.Clone()
	return &clone
}

func (user User) MattermostUsername() string {
	if user.MattermostUser == nil {
		return user.MattermostUserID
	}
	return user.MattermostUser.Username
}

func (user *User) AddEvent(event Event) {
	for _, existing := range user.Events {
		if existing == event.Event {
			return
		}
	}
	user.Events = append(user.Events, event.Event)
	eventsBy(byStartDate).Sort(user.Events)
}

func (user *User) overlapEvents(intervalStart, intervalEnd time.Time, remove bool) ([]store.Event, error) {
	var found, updated []store.Event
	for _, event := range user.Events {
		s, e, err := ParseDatePair(event.Start, event.End)
		if err != nil {
			return nil, err
		}

		// Find the overlap
		if s.Before(intervalStart) {
			s = intervalStart
		}
		if e.After(intervalEnd) {
			e = intervalEnd
		}

		if s.Before(e) {
			// Overlap
			found = append(found, event)
			if remove {
				continue
			}
		}

		updated = append(updated, event)
	}
	user.Events = updated
	return found, nil
}

func (api *api) loadOrMakeStoredUser(mattermostUserID string) (*User, bool, error) {
	storedUser, err := api.UserStore.LoadUser(mattermostUserID)
	var user *User
	if err == store.ErrNotFound {
		user, err = api.storeUserWelcomeNew(&User{
			User: store.NewUser(mattermostUserID),
		})
		return user, true, err
	}
	if err != nil {
		return nil, false, err
	}
	return &User{User: storedUser}, false, nil
}

// storeUserNotify checks if the user being stored is new, and welcomes the user.
// note that it can be used inside of filters, so it must not use filters itself,
//  nor assume that any runtime values have been filled.
func (api *api) storeUserWelcomeNew(orig *User) (*User, error) {
	user := orig.Clone()
	user.PluginVersion = api.Config.PluginVersion
	err := api.UserStore.StoreUser(user.User)
	if err != nil {
		return nil, err
	}

	if orig.PluginVersion == "" {
		api.messageWelcomeNewUser(user)
	}

	return user, nil
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

type UserMap map[string]*User

func (m UserMap) Clone(deep bool) UserMap {
	users := UserMap{}
	for id, user := range m {
		if deep {
			users[id] = user.Clone()
		} else {
			users[id] = user
		}
	}
	return users
}

func (m UserMap) IDMap() store.IDMap {
	ids := store.IDMap{}
	for id := range m {
		ids[id] = store.NotEmpty
	}
	return ids
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
