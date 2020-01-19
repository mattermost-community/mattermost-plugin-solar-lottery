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
		unmetNeeds:    unmetNeeds,
		unmetCapacity: unmetCapacity,
		orig:          errors.New("not ready"),
		shiftNumber:   shiftNumber,
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
			api.MarkdownUser(api.actingUser), api.MarkdownShift(rotation, shiftNumber))
		return nil, nil, nil
	}

	shift := shifts[0]
	added := addedUsers[0]
	logger.Infof("%s filled %s, added %s.",
		api.MarkdownUser(api.actingUser), api.MarkdownShift(rotation, shiftNumber), api.MarkdownUsersWithSkills(addedUsers[0]))
	return shift, added, nil
}
