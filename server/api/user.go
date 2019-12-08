// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type User interface {
	GetUser() (*store.User, error)
	JoinRotation(rotationName string, graceShifts int) error
	LeaveRotation(rotationName string) error
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

func withUser(api *api) error {
	if api.user != nil {
		return nil
	}

	user, err := api.UserStore.LoadUser(api.mattermostUserID)
	if err == store.ErrNotFound {
		api.user = store.NewUser(api.mattermostUserID)
		return nil
	}
	if err != nil {
		return err
	}
	api.user = user
	return nil
}
