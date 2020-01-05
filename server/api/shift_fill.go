// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrShiftMustBeOpen = errors.New("must be `open`")

func (api *api) IsShiftReady(rotation *Rotation, shiftNumber int) (shift *Shift, ready bool, whyNot string, err error) {
	shift, err = api.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, false, "", err
	}
	if shift.Status != store.ShiftStatusOpen {
		return nil, false, "", ErrShiftMustBeOpen
	}

	shiftUsers := rotation.ShiftUsers(shift)
	unmetNeeds := unmetNeeds(rotation.Needs, shiftUsers)
	unmetCapacity := 0
	if rotation.Size != 0 {
		unmetCapacity = rotation.Size - len(shift.MattermostUserIDs)
	}

	if len(unmetNeeds) == 0 && unmetCapacity <= 0 {
		return shift, true, "", nil
	}

	whyNot = autofillError{
		causeUnmetNeeds: unmetNeeds,
		causeCapacity:   unmetCapacity,
		orig:            errors.New("not ready"),
		shiftNumber:     shiftNumber,
	}.Error()

	return shift, false, whyNot, nil
}

func (api *api) FillShift(rotation *Rotation, shiftNumber int) (*Shift, UserMap, error) {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, nil, err
	}

	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.FillShifts",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	_, shifts, addedUsers, err := api.fillShifts(rotation, shiftNumber, 1, time.Time{}, logger)
	if err != nil {
		return nil, nil, err
	}

	if len(shifts) == 0 || len(addedUsers) == 0 {
		logger.Infof("%s tried to fill %v, nothing to do.",
			api.MarkdownUser(api.actingUser), MarkdownShift(rotation, shiftNumber))
		return nil, nil, nil
	}

	shift := shifts[0]
	added := addedUsers[0]
	logger.Infof("%s filled %s, added %s.",
		api.MarkdownUser(api.actingUser), MarkdownShift(rotation, shiftNumber), api.MarkdownUsersWithSkills(addedUsers[0]))
	return shift, added, nil
}

func (api *api) fillShifts(rotation *Rotation, startingShiftNumber, numShifts int, now time.Time, logger bot.Logger) ([]int, []*Shift, []UserMap, error) {
	// Guess' logs are too verbose - suppress
	prevLogger := api.Logger
	api.Logger = &bot.NilLogger{}
	shifts, err := api.Guess(rotation, startingShiftNumber, numShifts)
	api.Logger = prevLogger
	if err != nil {
		return nil, nil, nil, err
	}
	if len(shifts) != numShifts {
		return nil, nil, nil, errors.New("unreachable, must match")
	}

	var filledShiftNumbers []int
	var filledShifts []*Shift
	var addedUsers []UserMap

	appendShift := func(shiftNumber int, shift *Shift, added UserMap) {
		filledShiftNumbers = append(filledShiftNumbers, shiftNumber)
		filledShifts = append(filledShifts, shift)
		addedUsers = append(addedUsers, added)
	}

	shiftNumber := startingShiftNumber - 1
	for n := 0; n < numShifts; n++ {
		shiftNumber++

		loadedShift, err := api.OpenShift(rotation, shiftNumber)
		if err != nil && err != ErrShiftAlreadyExists {
			return nil, nil, nil, err
		}
		if loadedShift.Status != store.ShiftStatusOpen {
			appendShift(shiftNumber, loadedShift, nil)
			continue
		}
		if !loadedShift.Autopilot.Filled.IsZero() {
			appendShift(shiftNumber, loadedShift, nil)
			continue
		}

		before := rotation.ShiftUsers(loadedShift).Clone(false)

		// shifts coming from Guess are either loaded with their respective
		// status, or are Open. (in reality should always be Open).
		shift := shifts[n]
		added := UserMap{}
		for id, user := range rotation.ShiftUsers(shift) {
			if before[id] == nil {
				added[id] = user
			}
		}
		if len(added) == 0 {
			appendShift(shiftNumber, loadedShift, nil)
			continue
		}

		loadedShift.Autopilot.Filled = now

		_, err = api.joinShift(rotation, shiftNumber, loadedShift, added, true)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers, errors.WithMessagef(err, "failed to join autofilled users to %s", MarkdownShift(rotation, shiftNumber))
		}

		err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, loadedShift.Shift)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers, errors.WithMessagef(err, "failed to store autofilled %s", MarkdownShift(rotation, shiftNumber))
		}

		api.messageShiftJoined(added, rotation, shiftNumber, shift)
		appendShift(shiftNumber, shift, added)
	}

	return filledShiftNumbers, filledShifts, addedUsers, nil
}
