// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) JoinRotation(rotationName string, graceShifts int, mattermostUsername string) error {
	return api.userActRotation(rotationName, mattermostUsername,
		func(r *store.Rotation, user *store.User, username string, actingUser *store.User, actingUsername string) error {
			if len(r.MattermostUserIDs[user.MattermostUserID]) != 0 {
				return errors.Errorf("user %s (%s) is already in %s", username, user.MattermostUserID, r.Name)
			}

			shiftNumber, _ := ShiftNumber(r, time.Now())

			// A new person may be given some slack - setting LastShiftNumber in the
			// future guarantees they won't be selected until then.
			user.Rotations[rotationName] = shiftNumber + graceShifts

			err := api.UserStore.StoreUser(user)
			if err != nil {
				return err
			}

			if r.MattermostUserIDs == nil {
				r.MattermostUserIDs = store.UserIDList{}
			}
			r.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID

			message := "You have been added to the %q rotation"
			if actingUser != user {
				message += fmt.Sprintf(" by @%s", actingUsername)
			}
			message += "."
			api.Poster.DM(user.MattermostUserID, message, rotationName)

			api.Logger.With(bot.LogContext{
				"Username":       mattermostUsername,
				"ActingUsername": actingUsername,
				"RotationName":   rotationName,
				"LastShift":      shiftNumber + graceShifts,
			}).Debugf("%s has been added to %s.", mattermostUsername, rotationName)

			return nil
		})
}

func (api *api) LeaveRotation(rotationName, mattermostUsername string) error {
	return api.userActRotation(rotationName, mattermostUsername,
		func(r *store.Rotation, user *store.User, username string, actingUser *store.User, actingUsername string) error {
			if len(r.MattermostUserIDs) == 0 || len(r.MattermostUserIDs[user.MattermostUserID]) == 0 {
				return store.ErrNotFound
			}

			delete(user.Rotations, rotationName)
			err := api.UserStore.StoreUser(user)
			if err != nil {
				return err
			}

			delete(r.MattermostUserIDs, user.MattermostUserID)

			message := "You have been removed from the %q rotation"
			if actingUser != user {
				message += fmt.Sprintf(" by @%s", actingUsername)
			}
			message += "."
			api.Poster.DM(user.MattermostUserID, message, rotationName)

			api.Logger.With(bot.LogContext{
				"Username":       username,
				"ActingUsername": actingUsername,
				"RotationName":   rotationName,
			}).Debugf("%s has been removed from %s.", mattermostUsername, rotationName)

			return nil
		})
}

func (api *api) userActRotation(rotationName string, mattermostUsername string,
	updatef func(r *store.Rotation, user *store.User, username string, actingUser *store.User, actingUsername string) error) error {
	err := api.Filter(withUser, withRotations)
	if err != nil {
		return err
	}

	actingUser := api.user
	mmuser, err := api.PluginAPI.GetUser(actingUser.MattermostUserID)
	if err != nil {
		return errors.WithMessagef(err, "failed to load user `%s`", actingUser.MattermostUserID)
	}
	actingUsername := mmuser.Username

	user := api.user
	if mattermostUsername != "" {
		mmuser, err = api.PluginAPI.GetUserByUsername(mattermostUsername)
		if err != nil {
			return errors.WithMessagef(err, "failed to load user `%s`", mattermostUsername)
		}
		user, err = api.loadOrNewUser(mmuser.Id)
		if err != nil {
			return errors.WithMessagef(err, "failed to load user `%s`", mmuser.Id)
		}
	} else {
		mmuser, err = api.PluginAPI.GetUser(user.MattermostUserID)
		if err != nil {
			return errors.WithMessagef(err, "failed to load user `%s`", user.MattermostUserID)
		}
		mattermostUsername = mmuser.Username
	}

	r := api.rotations[rotationName]
	if r == nil {
		return store.ErrNotFound
	}

	err = updatef(r, user, mattermostUsername, actingUser, actingUsername)
	if err != nil {
		return err
	}

	return api.RotationsStore.StoreRotations(api.rotations)
}
