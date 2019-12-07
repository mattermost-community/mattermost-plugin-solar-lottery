// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-msoffice/server/utils/kvstore"
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
}

type Settings struct {
	EventSubscriptionID string
}

func (settings Settings) String() string {
	sub := "no subscription"
	if settings.EventSubscriptionID != "" {
		sub = "subscription ID: " + settings.EventSubscriptionID
	}
	return fmt.Sprintf(" - %s", sub)
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
	return kvstore.StoreJSON(s.userKV, user.MattermostUserID, user)
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	return s.userKV.Delete(mattermostUserID)
}
