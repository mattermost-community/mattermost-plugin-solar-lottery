// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type UserStore interface {
	LoadUser(mattermostUserId string) (*User, error)
	StoreUser(user *User) error
	DeleteUser(mattermostUserId string) error
}

const (
	StatusAvailable = ""
	StatusServing   = "serving"
	StatusBlocked   = "blocked"
)

type User struct {
	PluginVersion    string `json:",omitempty"`
	MattermostUserID string

	// Status is the user's current status
	Status string

	// Settings store the user's preferences.
	Settings Settings

	SkillLevels IntMap

	// Rotations is a map of all rotations (IDs) the user has joined. The
	// value is the last shift number served, for the rotation. When a user
	// joins a new rotation, their "last shift number" is set to the current
	// shift by default, offsetting it forward or backwards with a graceShifts
	// affects the new users likelihood of being selected for the next shift.
	// Setting it N shifts into the future guarantees that the new user will
	// not be selected until then.
	Rotations IntMap

	// Calendar is sorted by start date of the events
	Calendar []Event
}

const (
	EventTypeShift = "shift"
	EventTypeOther = "other"
)

type Event struct {
	Type string
	From string // time.RFC3339
	To   string // time.RFC3339
}

type Settings struct {
	Dummy bool
}

func NewUser(mattermostUserID string) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		SkillLevels:      IntMap{},
		Rotations:        IntMap{},
		Calendar:         []Event{},
	}
}

func (user *User) Clone() *User {
	clone := NewUser(user.MattermostUserID)
	clone.SkillLevels = user.SkillLevels.Clone()
	clone.Rotations = user.Rotations.Clone()
	clone.Calendar = append([]Event{}, user.Calendar...)
	return clone
}

func (s *pluginStore) LoadUser(mattermostUserId string) (*User, error) {
	user := NewUser("")
	err := kvstore.LoadJSON(s.userKV, mattermostUserId, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *pluginStore) StoreUser(user *User) error {
	err := kvstore.StoreJSON(s.userKV, user.MattermostUserID, user)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"User": user,
	}).Debugf("store: Stored user")
	return nil
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	err := s.userKV.Delete(mattermostUserID)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"MattermostUserID": mattermostUserID,
	}).Debugf("store: Deleted user")
	return nil
}
