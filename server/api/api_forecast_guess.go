// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) Guess(rotation *Rotation, startingShiftNumber int, numShifts int) ([]*Shift, error) {
	// Clone the rotation right away so we don't change the source.
	rotation = rotation.Clone(true)

	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.Guess",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"NumShifts":      numShifts,
		"ShiftNumber":    startingShiftNumber,
		"RotationID":     rotation.RotationID,
	})

	logger.Debugf("...running guess for\n%s", api.MarkdownRotationBullets(rotation))
	var shifts []*Shift
	for shiftNumber := startingShiftNumber; shiftNumber < startingShiftNumber+numShifts; shiftNumber++ {
		var shift *Shift
		shift, _, err := api.getShiftForGuess(rotation, shiftNumber)
		if err != nil {
			if err == store.ErrNotFound {
				shifts = append(shifts, nil)
				continue
			}
			return nil, err
		}

		if shift.Status == store.ShiftStatusOpen {
			err = api.autofillShift(rotation, shiftNumber, shift)
			if err != nil {
				return nil, err
			}
		}

		rotation.markShiftUsersEvents(shiftNumber, shift)
		rotation.markShiftUsersServed(shiftNumber, shift)
		shifts = append(shifts, shift)
	}

	logger.Debugf("Ran guess for %s", api.MarkdownRotation(rotation))
	return shifts, nil
}
