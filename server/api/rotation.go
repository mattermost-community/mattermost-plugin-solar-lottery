// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Rotation struct {
	*store.Rotation

	StartTime time.Time
	Users     UserMap
}

func (rotation *Rotation) init(api *api) error {
	start, err := time.Parse(DateFormat, rotation.Start)
	if err != nil {
		return err
	}
	rotation.StartTime = start

	if rotation.Users == nil {
		rotation.Users = UserMap{}
	}
	return nil
}

func (rotation *Rotation) ChangeNeed(skill string, level Level, newNeed store.Need) {
	for i, need := range rotation.Needs {
		if need.Skill == skill && need.Level == int(level) {
			rotation.Needs[i] = newNeed
			return
		}
	}
	rotation.Needs = append(rotation.Needs, newNeed)
}

func (rotation *Rotation) DeleteNeed(skill string, level Level) error {
	for i, need := range rotation.Needs {
		if need.Skill == skill && need.Level == int(level) {
			newNeeds := append([]store.Need{}, rotation.Needs[:i]...)
			if i+1 < len(rotation.Needs) {
				newNeeds = append(newNeeds, rotation.Needs[i+1:]...)
			}
			rotation.Needs = newNeeds
			return nil
		}
	}
	return errors.Errorf("%s is not found in rotation %s", MarkdownSkillLevel(skill, level), MarkdownRotation(rotation))
}

func (rotation *Rotation) ShiftNumberForTime(t time.Time) (int, error) {
	if t.Before(rotation.StartTime) {
		return -1, nil
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
			return -1, nil
		}
		if td < d {
			n--
		}
		return n, nil
	default:
		return -1, errors.Errorf("Invalid rotation period value %q", rotation.Period)
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
	if endShiftNumber == -1 || startShiftNumber == -1 || endShiftNumber > startShiftNumber {
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

func (api *api) ExpandRotation(rotation *Rotation) error {
	if !rotation.StartTime.IsZero() && len(rotation.Users) == len(rotation.MattermostUserIDs) {
		return nil
	}

	err := rotation.init(api)
	if err != nil {
		return err
	}

	users, err := api.LoadStoredUsers(rotation.MattermostUserIDs)
	if err != nil {
		return err
	}
	err = api.ExpandUserMap(users)
	if err != nil {
		return err
	}
	rotation.Users = users

	return nil
}

func withRotation(rotationID string) func(api *api) error {
	return func(api *api) error {
		return nil
	}
}

func withRotationExpanded(rotation *Rotation) func(api *api) error {
	return func(api *api) error {
		return api.ExpandRotation(rotation)
	}
}

func withRotationIsNotArchived(rotation *Rotation) func(api *api) error {
	return func(api *api) error {
		if rotation.IsArchived {
			return errors.Errorf("rotation %s is archived", MarkdownRotation(rotation))
		}
		return nil
	}
}
