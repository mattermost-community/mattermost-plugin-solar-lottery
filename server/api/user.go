// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type User interface {
	GetUser() (*store.User, error)
	JoinRotation(rotationName string, graceShifts int, mattermostUsername string) error
	LeaveRotation(rotationName, mattermostUsername string) error
	UpdateUserSkill(skill string, level int) (*store.User, error)
	DeleteUserSkill(skill string) (*store.User, error)
}

func (api *api) GetUser() (*store.User, error) {
	err := api.Filter(withUser)
	if err != nil {
		return nil, err
	}
	return api.user, nil
}

func (api *api) loadOrNewUser(mattermostUserID string) (*store.User, error) {
	user, err := api.UserStore.LoadUser(mattermostUserID)
	if err == store.ErrNotFound {
		return store.NewUser(mattermostUserID), nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func withUser(api *api) error {
	if api.user != nil {
		return nil
	}
	user, err := api.loadOrNewUser(api.mattermostUserID)
	if err != nil {
		return err
	}
	api.user = user
	return nil
}

func (api *api) loadUsers(ids store.UserIDList) (store.UserList, error) {
	users := store.UserList{}
	for id := range ids {
		u, err := api.UserStore.LoadUser(id)
		if err != nil {
			return nil, err
		}
		users[u.MattermostUserID] = u
	}
	return users, nil
}
