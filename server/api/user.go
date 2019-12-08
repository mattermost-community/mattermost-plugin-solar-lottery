// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"errors"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type User interface {
	JoinRotation(rotationName string) error
}

var ErrAlreadyInRotation = errors.New("user is already in rotation")

func (api *api) JoinRotation(rotationName string) error {
	err := api.Filter(withUser, withRotations)
	if err != nil {
		return nil
	}

	r := api.rotations[rotationName]
	if r == nil {
		return store.ErrNotFound
	}
	if len(r.MattermostUserIDs[api.user.MattermostUserID]) != 0 {
		return ErrAlreadyInRotation
	}

	prevUser := *api.user
	api.user.LastServedPeriod[rotationName] = -1
	err = api.UserStore.StoreUser(api.user)
	if err != nil {
		return err
	}

	r.MattermostUserIDs[api.user.MattermostUserID] = api.user.MattermostUserID
	err = api.RotationsStore.StoreRotations(api.rotations)
	if err != nil {
		_ = api.UserStore.StoreUser(&prevUser)
		return err
	}
	return nil
}

func withUser(api *api) error {
	if api.user != nil {
		return nil
	}

	user, err := api.UserStore.LoadUser(api.mattermostUserID)
	if err != nil {
		return err
	}

	api.user = user
	return nil
}
