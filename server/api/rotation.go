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

func (api *api) expandRotation(rotation *Rotation) error {
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
		return api.expandRotation(rotation)
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

func (rotation *Rotation) ChangeNeed(needName string, need store.Need) {
	if rotation.Needs == nil {
		rotation.Needs = map[string]store.Need{}
	}
	rotation.Needs[needName] = need
}

func (rotation *Rotation) DeleteNeed(needName string) error {
	_, ok := rotation.Needs[needName]
	if !ok {
		return errors.Errorf("%s is not found in rotation %s", needName, MarkdownRotation(rotation))
	}
	delete(rotation.Needs, needName)
	return nil
}

func (api *api) deleteUsersFromRotation(users UserMap, rotation *Rotation) error {
	for _, user := range users {
		_, ok := rotation.MattermostUserIDs[user.MattermostUserID]
		if !ok {
			return errors.Errorf("%s is not found in rotation %s", MarkdownUser(user), MarkdownRotation(rotation))
		}

		delete(user.Rotations, rotation.RotationID)
		_, err := api.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
		delete(rotation.MattermostUserIDs, user.MattermostUserID)
		api.messageLeftRotation(user, rotation)
		api.Logger.Debugf("removed %s from %s.", MarkdownUser(user), MarkdownRotation(rotation))
		return nil
	}

	return api.RotationStore.StoreRotation(rotation.Rotation)
}
