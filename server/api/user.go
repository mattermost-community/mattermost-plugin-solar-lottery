// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type User struct {
	*store.User

	// nil is assumed to be valid
	MattermostUser *model.User
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

func withActingUser(api *api) error {
	if api.actingUser != nil {
		return nil
	}
	user, _, err := api.loadOrMakeStoredUser(api.actingMattermostUserID)
	if err != nil {
		return err
	}
	api.actingUser = user
	return nil
}

func withActingUserExpanded(api *api) error {
	if api.actingUser != nil && api.actingUser.MattermostUser != nil {
		return nil
	}
	err := withActingUser(api)
	if err != nil {
		return err
	}
	return api.ExpandUser(api.actingUser)
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
func (api *api) storeUserWelcomeNew(u *User) (*User, error) {
	user := u.Clone()
	user.PluginVersion = api.Config.PluginVersion
	err := api.UserStore.StoreUser(user.User)
	if err != nil {
		return nil, err
	}

	api.messageWelcomeNewUser(user)
	return user, nil
}
