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

	ShiftUsers := rotation.ShiftUsers(shift)
	unsatisfiedNeeds := api.unsatisfiedNeeds(rotation.Needs, ShiftUsers)
	unsatisfiedCapacity := 0
	if rotation.Size != 0 {
		unsatisfiedCapacity = rotation.Size - len(shift.MattermostUserIDs)
	}

	if len(unsatisfiedNeeds) == 0 && unsatisfiedCapacity <= 0 {
		return shift, true, "", nil
	}

	whyNot = autofillError{
		unsatisfiedNeeds:    unsatisfiedNeeds,
		unsatisfiedCapacity: unsatisfiedCapacity,
		orig:                errors.New("not ready"),
		shiftNumber:         shiftNumber,
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

	shifts, addedUsers, err := api.fillShifts(rotation, shiftNumber, 1, false, time.Time{}, logger)
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

func (api *api) fillShifts(rotation *Rotation, shiftNumber, numShifts int, autofill bool, now time.Time, logger bot.Logger) ([]*Shift, []UserMap, error) {
	shifts, err := api.Guess(rotation, shiftNumber, numShifts, true)
	if err != nil {
		return nil, nil, err
	}
	if len(shifts) != numShifts {
		return nil, nil, errors.New("unreachable, must match")
	}

	var filledShifts []*Shift
	var addedUsers []UserMap
	for n := shiftNumber; n < shiftNumber+numShifts; n++ {
		loadedShift, err := api.OpenShift(rotation, n)
		if err != nil && err != ErrShiftAlreadyExists {
			return nil, nil, err
		}
		if loadedShift.Status != store.ShiftStatusOpen {
			logger.Debugf("<><> ignored %s, status %s", MarkdownShift(rotation, shiftNumber), loadedShift.Status)
			continue
		}
		if autofill && !loadedShift.Autopilot.Filled.IsZero() {
			logger.Debugf("<><> already filled on %v", loadedShift.Autopilot.Filled)
			continue
		}

		before := rotation.ShiftUsers(loadedShift).Clone(false)

		shift := shifts[n]
		added := UserMap{}
		for id, user := range rotation.ShiftUsers(shift) {
			if before[id] == nil {
				added[id] = user
			}
		}
		if len(added) == 0 {
			logger.Debugf("<><> ignored %s, no users added", MarkdownShift(rotation, shiftNumber))
			continue
		}

		if autofill {
			shift.Autopilot.Filled = now
		}

		err = api.ShiftStore.StoreShift(rotation.RotationID, n, shift.Shift)
		if err != nil {
			return nil, nil, errors.WithMessagef(err, "failed to store autofilled %s", MarkdownShift(rotation, shiftNumber))
		}

		api.messageShiftJoined(added, rotation, shiftNumber, shift)

		filledShifts = append(filledShifts, shift)
		addedUsers = append(addedUsers, added)
	}

	return filledShifts, addedUsers, nil
}
