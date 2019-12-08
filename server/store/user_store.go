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
	PluginVersion    string
	MattermostUserID string
	Settings         Settings `json:"mattermostSettings,omitempty"`
	LastServedPeriod map[string]int
	Unavailables     []Unavailable
}

type Unavailable struct {
	From    time.Time
	To      time.Time
	Comment string
}

type Settings struct {
	Dummy bool
}

func (settings Settings) String() string {
	return "settings <><>"
}

func (s *pluginStore) LoadUser(mattermostUserId string) (*User, error) {
	user := User{}
	err := kvstore.LoadJSON(s.userKV, mattermostUserId, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
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
