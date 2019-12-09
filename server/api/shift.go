// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

const Week = time.Hour * 24 * 7
const DateFormat = "2006-01-02"

func ShiftNumber(r *store.Rotation, t time.Time) (int, error) {
	start, err := time.Parse(DateFormat, r.Start)
	if err != nil {
		return 0, err
	}
	if time.Now().Before(start) {
		return 0, errors.Errorf("Time %v is before rotation start %v", t, start)
	}

	switch r.Period {
	case "1w", "w":
		return int(t.Sub(start) / Week), nil
	case "2w":
		return int(t.Sub(start) / (2 * Week)), nil
	case "1m", "m":
		y, m, d := start.Date()
		ty, tm, td := t.Date()
		n := (ty*12 + int(tm)) - (y*12 - int(m))
		if td >= d {
			n++
		}
		return n, nil
	default:
		return 0, errors.Errorf("Invalid rotation period value %q", r.Period)
	}
	return 0, nil
}

func ShiftDates(r *store.Rotation, shiftNumber int) (time.Time, time.Time, error) {
	rstart, err := time.Parse(DateFormat, r.Start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	var begin, end time.Time
	switch r.Period {
	case "1w", "w":
		begin = rstart.Add(time.Duration(shiftNumber) * Week)
		end = begin.Add(Week)

	case "2w":
		begin = rstart.Add(time.Duration(shiftNumber) * 2 * Week)
		end = begin.Add(2 * Week)

	case "1m", "m":
		y, month, d := rstart.Date()
		m := int(month-1) + shiftNumber
		year := y + m/12
		month = time.Month((m % 12) + 1)
		begin = time.Date(year, month, d, 0, 0, 0, 0, rstart.Location())
		m++
		year = y + m/12
		month = time.Month((m % 12) + 1)
		end = time.Date(year, month, d, 0, 0, 0, 0, rstart.Location())

	default:
		return time.Time{}, time.Time{}, errors.Errorf("Invalid rotation period value %q", r.Period)
	}
	return begin, end, nil
}
