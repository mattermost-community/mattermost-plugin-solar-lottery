// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"errors"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

var ErrAlreadyInRotation = errors.New("user is already in rotation")

func (api *api) JoinRotation(rotationName string, graceShifts int) error {
	err := api.Filter(withUser, withRotations)
	if err != nil {
		return err
	}

	r := api.rotations[rotationName]
	if r == nil {
		return store.ErrNotFound
	}
	if len(r.MattermostUserIDs[api.user.MattermostUserID]) != 0 {
		return ErrAlreadyInRotation
	}

	shiftNumber, _ := ShiftNumber(r, time.Now())

	// A new person may be given some slack - setting LastShiftNumber in the
	// future guarantees they won't be selected until then.
	api.user.Rotations[rotationName] = shiftNumber + graceShifts

	err = api.UserStore.StoreUser(api.user)
	if err != nil {
		return err
	}

	if r.MattermostUserIDs == nil {
		r.MattermostUserIDs = store.UserIDList{}
	}
	r.MattermostUserIDs[api.user.MattermostUserID] = api.user.MattermostUserID

	return api.RotationsStore.StoreRotations(api.rotations)
}

func (api *api) LeaveRotation(rotationName string) error {
	err := api.Filter(withUser, withRotations)
	if err != nil {
		return nil
	}

	r := api.rotations[rotationName]
	if r == nil {
		return store.ErrNotFound
	}
	if len(r.MattermostUserIDs) == 0 || len(r.MattermostUserIDs[api.user.MattermostUserID]) == 0 {
		return store.ErrNotFound
	}

	user := api.user
	delete(user.Rotations, rotationName)
	err = api.UserStore.StoreUser(api.user)
	if err != nil {
		return err
	}

	delete(r.MattermostUserIDs, api.user.MattermostUserID)
	return api.RotationsStore.StoreRotations(api.rotations)
}
