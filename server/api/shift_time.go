// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"
)

func (rotation *Rotation) ShiftNumberForTime(t time.Time) (int, error) {
	if t.Before(rotation.StartTime) {
		return 0, errors.Errorf("Time %v is before rotation start %v", t, rotation.StartTime)
	}

	switch rotation.Period {
	case EveryWeek:
		return int(t.Sub(rotation.StartTime) / WeekDuration), nil
	case EveryTwoWeeks:
		return int(t.Sub(rotation.StartTime) / (2 * WeekDuration)), nil
	case EveryMonth:
		y, m, d := rotation.StartTime.Date()
		ty, tm, td := t.Date()
		n := (ty*12 + int(tm)) - (y*12 + int(m))
		if n <= 0 {
			return 0, nil
		}
		if td < d {
			n--
		}
		return n, nil
	default:
		return 0, errors.Errorf("Invalid rotation period value %q", rotation.Period)
	}
}

func (rotation *Rotation) ShiftsForDates(startDate, endDate string) (first, numShifts int, err error) {
	start, err := time.Parse(DateFormat, startDate)
	if err != nil {
		return 0, 0, err
	}
	startShiftNumber, err := rotation.ShiftNumberForTime(start)
	if err != nil {
		return 0, 0, err
	}
	end, err := time.Parse(DateFormat, endDate)
	if err != nil {
		return 0, 0, err
	}
	endShiftNumber, err := rotation.ShiftNumberForTime(end)
	if err != nil {
		return 0, 0, err
	}
	if endShiftNumber > startShiftNumber {
		return 0, 0, errors.Errorf("invalid date range: from %s to %s", startDate, endDate)
	}

	return startShiftNumber, endShiftNumber - startShiftNumber + 1, nil
}

func (rotation *Rotation) ShiftDatesForNumber(shiftNumber int) (time.Time, time.Time, error) {
	var begin, end time.Time
	switch rotation.Period {
	case EveryWeek:
		begin = rotation.StartTime.Add(time.Duration(shiftNumber) * WeekDuration)
		end = begin.Add(WeekDuration)

	case EveryTwoWeeks:
		begin = rotation.StartTime.Add(time.Duration(shiftNumber) * 2 * WeekDuration)
		end = begin.Add(2 * WeekDuration)

	case EveryMonth:
		y, month, d := rotation.StartTime.Date()
		m := int(month-1) + shiftNumber
		year := y + m/12
		month = time.Month((m % 12) + 1)
		begin = time.Date(year, month, d, 0, 0, 0, 0, rotation.StartTime.Location())
		m++
		year = y + m/12
		month = time.Month((m % 12) + 1)
		end = time.Date(year, month, d, 0, 0, 0, 0, rotation.StartTime.Location())

	default:
		return time.Time{}, time.Time{}, errors.Errorf("Invalid rotation period value %q", rotation.Period)
	}
	return begin, end, nil
}
