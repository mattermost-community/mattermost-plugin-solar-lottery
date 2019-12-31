// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

const WeekDuration = time.Hour * 24 * 7
const DateFormat = "2006-01-02"

type Shift struct {
	*store.Shift
	StartTime time.Time
	EndTime   time.Time
	Users     UserMap
}

func (api *api) ExpandShift(shift *Shift) error {
	if !shift.StartTime.IsZero() && len(shift.Users) == len(shift.MattermostUserIDs) {
		return nil
	}

	err := shift.init()
	if err != nil {
		return err
	}

	users, err := api.LoadStoredUsers(shift.MattermostUserIDs)
	if err != nil {
		return err
	}
	err = api.ExpandUserMap(users)
	if err != nil {
		return err
	}
	shift.Users = users

	return nil
}

func (api *api) loadOrMakeOneShift(rotation *Rotation, shiftNumber int, autofill bool) (*Shift, bool, error) {
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
		err = api.ExpandShift(shift)

	case store.ErrNotFound:
		if !autofill {
			return nil, false, err
		}
		shift, err = rotation.makeShift(shiftNumber, nil)
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

	err = api.ExpandShift(shift)
	if err != nil {
		return nil, false, err
	}

	return shift, created, nil
}

func (rotation *Rotation) makeShift(shiftNumber int, users UserMap) (*Shift, error) {
	start, end, err := rotation.ShiftDatesForNumber(shiftNumber)
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = UserMap{}
	}
	return &Shift{
		Shift:     store.NewShift(start.Format(DateFormat), end.Format(DateFormat), users.IDMap()),
		StartTime: start,
		EndTime:   end,
		Users:     users,
	}, nil
}

func (api *api) loadShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := rotation.makeShift(shiftNumber, nil)
	if err != nil {
		return nil, err
	}
	s, err := api.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	if err != nil {
		return nil, err
	}
	shift.Shift = s
	err = api.ExpandShift(shift)
	if err != nil {
		return nil, err
	}
	return shift, nil
}

func (shift *Shift) init() error {
	start, err := time.Parse(DateFormat, shift.Start)
	if err != nil {
		return err
	}
	end, err := time.Parse(DateFormat, shift.End)
	if err != nil {
		return err
	}
	shift.StartTime = start
	shift.EndTime = end
	shift.Users = UserMap{}
	return nil
}
