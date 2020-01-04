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

	// Map of last shift served per rotation, used to calculate weight
	LastServed IntMap

	// Events is sorted by start date of the events.
	Events []Event
}

const (
	EventTypeShift       = "shift"
	EventTypeUnavailable = "unavavailable"
)

type Event struct {
	Type  string
	Start string // time.RFC3339
	End   string // time.RFC3339

	// Rotation and ShiftNumber identify the shift for Shift and Padding event
	// types.
	RotationID  string
	ShiftNumber int
}

type Settings struct {
	Dummy bool
}

func NewUser(mattermostUserID string) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		SkillLevels:      IntMap{},
		LastServed:       IntMap{},
		Events:           []Event{},
	}
}

func (user *User) Clone() *User {
	clone := NewUser(user.MattermostUserID)
	clone.SkillLevels = user.SkillLevels.Clone()
	clone.LastServed = user.LastServed.Clone()
	clone.Events = append([]Event{}, user.Events...)
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
