// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

func (api *api) Guess(rotation *Rotation, startingShiftNumber int, numShifts int, autofill bool) ([]*Shift, error) {
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
		"Autofill":       autofill,
		"ShiftNumber":    startingShiftNumber,
		"RotationID":     rotation.RotationID,
	})

	var shifts []*Shift
	for shiftNumber := startingShiftNumber; shiftNumber < startingShiftNumber+numShifts; shiftNumber++ {
		var shift *Shift
		shift, _, err := api.getShiftForGuess(rotation, shiftNumber, autofill)
		if err != nil {
			if !autofill && err == store.ErrNotFound {
				shifts = append(shifts, nil)
				continue
			}
			return nil, err
		}

		if autofill && shift.Status == store.ShiftStatusOpen {
			err = api.autofillShift(rotation, shiftNumber, shift, autofill)
			if err != nil {
				return nil, err
			}
		}

		rotation.markShiftUsersEvents(shiftNumber, shift)
		rotation.markShiftUsersServed(shiftNumber, shift)
		shifts = append(shifts, shift)
	}

	logger.Debugf("Ran guess for %s", MarkdownRotation(rotation))
	return shifts, nil
}

// Returns an un-expanded shift - will be populated with Users from rotation
func (api *api) getShiftForGuess(rotation *Rotation, shiftNumber int, autofill bool) (*Shift, bool, error) {
	start, end, err := rotation.ShiftDatesForNumber(shiftNumber)
	if err != nil {
		return nil, false, err
	}

	var shift *Shift
	created := false
	storedShift, err := api.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	switch err {
	case nil:
		shift = &Shift{
			Shift: storedShift,
		}

	case store.ErrNotFound:
		if !autofill {
			return nil, false, err
		}
		shift, err = rotation.makeShift(shiftNumber)
		if err != nil {
			return nil, false, err
		}
		created = true

	default:
		return nil, false, err
	}

	if shift.Start != start.Format(DateFormat) || shift.End != end.Format(DateFormat) {
		return nil, false, errors.Errorf("loaded shift has wrong dates %v-%v, expected %v-%v",
			shift.Start, shift.End, start, end)
	}

	return shift, created, nil
}
