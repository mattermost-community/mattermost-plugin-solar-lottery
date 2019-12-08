// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type User interface {
	AddUserToRotation(user *store.User, rotation *store.Rotation) error
}

var ErrAlreadyInRotation = errors.New("user is already in rotation")

func (api *api) AddUserToRotation(user *store.User, rotation *store.Rotation) error {
	if len(rotation.MattermostUserIDs[user.MattermostUserID]) != 0 {
		return ErrAlreadyInRotation
	}
	prevUser := *user
	user.LastServedPeriod[rotation.Name] = -1
	err := api.UserStore.StoreUser(user)
	if err != nil {
		return err
	}
	rotation.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID
	err = api.RotationStore.StoreRotation(rotation)
	if err != nil {
		_ = api.UserStore.StoreUser(&prevUser)
		return err
	}
	return nil
}
