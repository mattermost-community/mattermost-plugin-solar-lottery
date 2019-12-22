// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

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

func (api *api) expandShift(shift *Shift) error {
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