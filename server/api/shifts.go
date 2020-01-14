// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
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

// Returns an un-expanded shift - will be populated with Users from rotation
func (api *api) getShiftForGuess(rotation *Rotation, shiftNumber int) (*Shift, bool, error) {
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

	case store.ErrNotFound:
		shift, err = rotation.makeShift(shiftNumber)
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

	return shift, created, nil
}

func (api *api) startShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}
	if shift.Status == store.ShiftStatusStarted {
		return shift, errors.New("already started")
	}
	if shift.Status != store.ShiftStatusOpen {
		return nil, errors.Errorf("can't start a shift which is %s, must be open", shift.Status)
	}

	shift.Status = store.ShiftStatusStarted

	for _, user := range rotation.ShiftUsers(shift) {
		rotation.markShiftUserServed(user, shiftNumber, shift)
		_, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return nil, err
		}
	}

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	api.messageShiftStarted(rotation, shiftNumber, shift)
	return shift, nil
}

func (api *api) finishShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}
	if shift.Status == store.ShiftStatusFinished {
		return shift, nil
	}
	if shift.Status != store.ShiftStatusStarted {
		return nil, errors.Errorf("can't finish a shift which is %s, must be started", shift.Status)
	}

	shift.Status = store.ShiftStatusFinished
	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	api.messageShiftFinished(rotation, shiftNumber, shift)

	return shift, nil
}

func (api *api) joinShift(rotation *Rotation, shiftNumber int, shift *Shift, users UserMap, persist bool) (UserMap, error) {
	if shift.Status != store.ShiftStatusOpen {
		return nil, errors.Errorf("can't join a shift with status %s, must be Open", shift.Status)
	}

	joined := UserMap{}
	for _, user := range users {
		if shift.Shift.MattermostUserIDs[user.MattermostUserID] != "" {
			continue
		}
		if len(shift.MattermostUserIDs) >= rotation.Size {
			return nil, errors.Errorf("rotation size %v exceeded", rotation.Size)
		}
		shift.Shift.MattermostUserIDs[user.MattermostUserID] = store.NotEmpty
		joined[user.MattermostUserID] = user
	}

	err := api.addEventToUsers(joined, NewShiftEvent(rotation, shiftNumber, shift), persist)
	if err != nil {
		return nil, err
	}

	return joined, nil
}

func (api *api) fillShifts(rotation *Rotation, startingShiftNumber, numShifts int, now time.Time, logger bot.Logger) ([]int, []*Shift, []UserMap, error) {
	// Guess' logs are too verbose - suppress
	prevLogger := api.Logger
	api.Logger = &bot.NilLogger{}
	shifts, err := api.Guess(rotation, startingShiftNumber, numShifts)
	api.Logger = prevLogger
	if err != nil {
		return nil, nil, nil, err
	}
	if len(shifts) != numShifts {
		return nil, nil, nil, errors.New("unreachable, must match")
	}

	var filledShiftNumbers []int
	var filledShifts []*Shift
	var addedUsers []UserMap

	appendShift := func(shiftNumber int, shift *Shift, added UserMap) {
		filledShiftNumbers = append(filledShiftNumbers, shiftNumber)
		filledShifts = append(filledShifts, shift)
		addedUsers = append(addedUsers, added)
	}

	shiftNumber := startingShiftNumber - 1
	for n := 0; n < numShifts; n++ {
		shiftNumber++

		loadedShift, err := api.OpenShift(rotation, shiftNumber)
		if err != nil && err != ErrShiftAlreadyExists {
			return nil, nil, nil, err
		}
		if loadedShift.Status != store.ShiftStatusOpen {
			appendShift(shiftNumber, loadedShift, nil)
			continue
		}
		if !loadedShift.Autopilot.Filled.IsZero() {
			appendShift(shiftNumber, loadedShift, nil)
			continue
		}

		before := rotation.ShiftUsers(loadedShift).Clone(false)

		// shifts coming from Guess are either loaded with their respective
		// status, or are Open. (in reality should always be Open).
		shift := shifts[n]
		added := UserMap{}
		for id, user := range rotation.ShiftUsers(shift) {
			if before[id] == nil {
				added[id] = user
			}
		}
		if len(added) == 0 {
			appendShift(shiftNumber, loadedShift, nil)
			continue
		}

		loadedShift.Autopilot.Filled = now

		_, err = api.joinShift(rotation, shiftNumber, loadedShift, added, true)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers, errors.WithMessagef(err, "failed to join autofilled users to %s", api.MarkdownShift(rotation, shiftNumber))
		}

		err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, loadedShift.Shift)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers, errors.WithMessagef(err, "failed to store autofilled %s", api.MarkdownShift(rotation, shiftNumber))
		}

		api.messageShiftJoined(added, rotation, shiftNumber, shift)
		appendShift(shiftNumber, shift, added)
	}

	return filledShiftNumbers, filledShifts, addedUsers, nil
}
