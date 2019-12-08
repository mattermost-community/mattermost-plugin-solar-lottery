// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type UserStore interface {
	LoadUser(mattermostUserId string) (*User, error)
	StoreUser(user *User) error
	DeleteUser(mattermostUserId string) error
}

type User struct {
	PluginVersion    string `json:",omitempty"`
	MattermostUserID string

	// Settings store the user's preferences.
	Settings Settings `json:"mattermostSettings,omitempty"`

	// Joined is a map of all subscription (names) the user has joined. The
	// value is the last shift number served, for the rotation. When a user
	// joins a new rotation, their "last shift number" is set to the current
	// period by default, offsetting it forward or backwards with a graceShifts
	// affects the new users likelihood of being selected for the next shift.
	// Setting it N shifts into the future guarantees that the new user will
	// not be selected until then.
	Joined map[string]int `json:",omitempty"`

	// Unavailables stores the times of user unavailability, applies to all
	// rotations the user is in.
	Unavailables []Unavailable `json:",omitempty"`
}

type Unavailable struct {
	From    time.Time
	To      time.Time
	Comment string
}

type Settings struct {
	Dummy bool
}

func NewUser(mattermostUserID string) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		Joined:           map[string]int{},
		Unavailables:     []Unavailable{},
	}
}

func (s *pluginStore) LoadUser(mattermostUserId string) (*User, error) {
	user := NewUser(mattermostUserId)
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
	}).Debugf("Stored user")
	return nil
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	err := s.userKV.Delete(mattermostUserID)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"MattermostUserID": mattermostUserID,
	}).Debugf("Deleted user")
	return nil
}
