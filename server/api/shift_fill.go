// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrShiftMustBeOpen = errors.New("must be `open`")

func (api *api) FillShift(rotation *Rotation, shiftNumber int, autofill bool) (shift *Shift, ready bool, whyNot string, before, added UserMap, err error) {
	err = api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, false, "", nil, nil, err
	}
	shift, err = api.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, false, "", nil, nil, err
	}
	if shift.Status != store.ShiftStatusOpen {
		return nil, false, "", nil, nil, ErrShiftMustBeOpen
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.FillShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	unfulfilledNeeds := api.unsatisfiedNeeds(rotation.Needs, shift.Users)
	unfulfilledCapacity := 0
	if rotation.Size != 0 {
		unfulfilledCapacity = rotation.Size - len(shift.Users)
	}

	if len(unfulfilledNeeds) == 0 || unfulfilledCapacity <= 0 {
		return shift, true, "", shift.Users, nil, nil
	}
	if !autofill {
		return shift,
			false,
			autofillError{
				unfulfilledNeeds:    api.unsatisfiedNeeds(rotation.Needs, shift.Users),
				unfulfilledCapacity: rotation.Size - len(shift.Users),
				orig:                errors.New("not ready"),
				shiftNumber:         shiftNumber,
			}.Error(),
			shift.Users, nil, nil
	}

	shifts, err := api.Guess(rotation, shiftNumber, 1, autofill)
	if err != nil {
		return nil, false, "", nil, nil, errors.WithMessagef(err, "failed to autofill %s", MarkdownShift(rotation, shiftNumber, shift))
	}

	before = shift.Users.Clone(false)
	shift = shifts[0]
	added = UserMap{}
	for id, user := range shift.Users {
		if before[id] == nil {
			added[id] = user
		}
	}

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, false, "", nil, nil, errors.WithMessagef(err, "failed to store autofilled %s", MarkdownShift(rotation, shiftNumber, shift))
	}

	api.messageShiftVolunteers(added, rotation, shiftNumber, shift)
	logger.Infof("%s filled %s, added %s.",
		MarkdownUser(api.actingUser), MarkdownShift(rotation, shiftNumber, shift), MarkdownUserMapWithSkills(added))

	return shift, true, "", before, added, nil
}
