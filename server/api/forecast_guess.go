// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (api *api) Guess(rotation *Rotation, startingShiftNumber int, numShifts int, autofill bool) ([]*Shift, error) {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	// logger := api.Logger.Timed().With(bot.LogContext{
	// 	"Location":       "api.Guess",
	// 	"ActingUsername": api.actingUser.MattermostUsername(),
	// 	"NumShifts":      numShifts,
	// 	"Autofill":       autofill,
	// 	"ShiftNumber":    startingShiftNumber,
	// 	"RotationID":     rotation.RotationID,
	// })

	prevUsers := rotation.Users
	rotation.Users = rotation.Users.Clone(true)
	defer func() {
		rotation.Users = prevUsers
	}()

	// Need to cache all users between iterations, to mock their last served dates.
	cachedUsers := UserMap{}
	for mattermostUserID, user := range rotation.Users {
		cachedUsers[mattermostUserID] = user
	}

	var shifts []*Shift
	for shiftNumber := startingShiftNumber; shiftNumber < startingShiftNumber+numShifts; shiftNumber++ {
		var shift *Shift
		shift, _, err := api.loadOrMakeOneShift(rotation, shiftNumber, autofill)
		if err != nil {
			if !autofill && err == store.ErrNotFound {
				shifts = append(shifts, nil)
				continue
			}
			return nil, err
		}

		// If the Shift was loaded, it might have brought some new users with
		// it. Replace them with cached Users as appropriate, they are more up-
		// to-date.
		for k := range shift.Users {
			if cachedUsers[k] != nil {
				shift.Users[k] = cachedUsers[k]
			}
		}

		if autofill && (shift.Status == "" || shift.Status == store.ShiftStatusOpen) {
			err = api.autofillShift(rotation, shiftNumber, shift, autofill)
			if err != nil {
				return nil, err
			}
		}

		// Update shift's users' last served counter, and update the cache in
		// case they were not there
		for k, u := range shift.Users {
			u.NextRotationShift[rotation.RotationID] = shiftNumber + 1 + rotation.Grace
			cachedUsers[k] = u
		}

		shifts = append(shifts, shift)
	}

	// logger.Debugf("Ran guess for %s", MarkdownRotation(rotation))
	return shifts, nil
}
