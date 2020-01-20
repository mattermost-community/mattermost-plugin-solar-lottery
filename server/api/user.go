// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"strings"
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

func (api *api) ExpandUser(user *User) error {
	if user.MattermostUser != nil {
		return nil
	}
	mattermostUser, err := api.PluginAPI.GetMattermostUser(user.MattermostUserID)
	if err != nil {
		return err
	}
	user.MattermostUser = mattermostUser
	return nil
}

func (user *User) WithLastServed(rotationID string, shiftNumber int) *User {
	newUser := user.Clone()
	newUser.LastServed[rotationID] = shiftNumber
	return newUser
}

func (user *User) WithSkills(skillsLevels store.IntMap) *User {
	newUser := user.Clone()
	if newUser.SkillLevels != nil {
		newUser.SkillLevels = store.IntMap{}
	}
	for s, l := range skillsLevels {
		newUser.SkillLevels[s] = l
	}
	return newUser
}

func (user *User) String() string {
	if user.MattermostUser != nil {
		return fmt.Sprintf("@%s", user.MattermostUser.Username)
	} else {
		return fmt.Sprintf("%q", user.MattermostUserID)
	}
}

func (user *User) Markdown() string {
	if user.MattermostUser != nil {
		return fmt.Sprintf("@%s", user.MattermostUser.Username)
	} else {
		return fmt.Sprintf("userID `%s`", user.MattermostUserID)
	}
}

func (user *User) MarkdownWithSkills() string {
	return fmt.Sprintf("%s %s", user.Markdown(), user.MarkdownSkills())
}

func (user *User) MarkdownSkills() string {
	skills := []string{}
	for s, l := range user.SkillLevels {
		skills = append(skills, MarkdownSkillLevel(s, Level(l)))
	}

	if len(skills) == 0 {
		return "(kook)"
	}
	ss := strings.Join(skills, ", ")
	return fmt.Sprintf("(%s)", ss)
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

func (user *User) OverlapEvents(intervalStart, intervalEnd time.Time, remove bool) ([]store.Event, error) {
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
		api.Logger.Debugf("nothing to do for user %s, already has skill %s (%v)", user.Markdown(), skillName, level)
		return user, nil
	}

	if level == 0 {
		_, ok := user.SkillLevels[skillName]
		if !ok {
			return nil, errors.Errorf("%s does not have skill %s", user.Markdown(), skillName)
		}
		delete(user.SkillLevels, skillName)
	} else {
		user.SkillLevels[skillName] = int(level)
	}

	user, err := api.storeUserWelcomeNew(user)
	if err != nil {
		return nil, err
	}
	api.Logger.Debugf("%s (%v) skill updated user %s", skillName, level, user.Markdown())
	return user, nil
}

