// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

const DayDuration = time.Hour * 24
const WeekDuration = DayDuration * 7
const DateFormat = "2006-01-02"

type Shift struct {
	*store.Shift
	StartTime time.Time
	EndTime   time.Time
}

func (shift *Shift) Clone(deep bool) *Shift {
	newShift := *shift
	if deep {
		newShift.Shift = &(*shift.Shift)
	}
	return &newShift
}

func (rotation *Rotation) makeShift(shiftNumber int) (*Shift, error) {
	start, end, err := rotation.ShiftDatesForNumber(shiftNumber)
	if err != nil {
		return nil, err
	}
	return &Shift{
		Shift:     store.NewShift(start.Format(DateFormat), end.Format(DateFormat), nil),
		StartTime: start,
		EndTime:   end,
	}, nil
}

func (api *api) loadShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := rotation.makeShift(shiftNumber)
	if err != nil {
		return nil, err
	}
	s, err := api.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	if err != nil {
		return nil, err
	}
	shift.Shift = s

	shift.StartTime, shift.EndTime, err = ParseDatePair(s.Start, s.End)
	if err != nil {
		return nil, err
	}

	return shift, nil
}
