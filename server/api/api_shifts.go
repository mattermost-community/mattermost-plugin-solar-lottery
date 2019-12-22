// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Shifts interface {
	CommitShift(rotation *Rotation, shiftNumber int) error
	ListShifts(rotation *Rotation, startDate string, numShifts int) ([]*Shift, error)
	OpenShift(*Rotation, int) (*Shift, error)
	StartShift(rotation *Rotation, shiftNumber int) error
	FinishShift(rotation *Rotation, shiftNumber int) error

	DebugDeleteShift(rotation *Rotation, shiftNumber int) error
}

var ErrShiftAlreadyExists = errors.New("shift already exists")

func (api *api) CommitShift(rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.CommitShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return err
	}

	shift.ShiftStatus = store.ShiftStatusCommitted

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return err
	}

	api.messageShiftCommitted(rotation, shiftNumber, shift)
	logger.Infof("%s committed shift %s in %s.", MarkdownUser(api.actingUser), MarkdownShift(shiftNumber, shift), MarkdownRotation(rotation))
	return nil
}

func (api *api) StartShift(rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.StartShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return err
	}

	shift.ShiftStatus = store.ShiftStatusStarted

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return err
	}

	api.messageShiftStarted(rotation, shiftNumber, shift)
	logger.Infof("%s started shift %s in %s.", MarkdownUser(api.actingUser), MarkdownShift(shiftNumber, shift), MarkdownRotation(rotation))
	return nil
}

func (api *api) DebugDeleteShift(rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.DebugDeleteShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	err = api.ShiftStore.DeleteShift(rotation.RotationID, shiftNumber)
	if err != nil {
		return err
	}

	logger.Infof("%s deleted shift %v in %s.", MarkdownUser(api.actingUser), shiftNumber, MarkdownRotation(rotation))
	return nil
}

func (api *api) FinishShift(rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.FinishShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return err
	}

	shift.ShiftStatus = store.ShiftStatusFinished

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return err
	}

	api.messageShiftFinished(rotation, shiftNumber, shift)
	logger.Infof("%s finished shift %s in %s.", MarkdownUser(api.actingUser), MarkdownShift(shiftNumber, shift), MarkdownRotation(rotation))
	return nil
}

func (api *api) ListShifts(rotation *Rotation, startDate string, numShifts int) ([]*Shift, error) {
	err := api.Filter(
	// withActingUser,
	)
	if err != nil {
		return nil, err
	}

	starting, err := time.Parse(DateFormat, startDate)
	if err != nil {
		return nil, err
	}
	n, err := rotation.ShiftNumberForTime(starting)
	if err != nil {
		return nil, err
	}

	shifts := []*Shift{}
	for i := n; i < n+numShifts; i++ {
		var shift *Shift
		shift, err = api.loadShift(rotation, i)
		if err != nil {
			if err != store.ErrNotFound {
				return nil, err
			}
			continue
		}

		shifts = append(shifts, shift)
	}

	return shifts, nil
}

func (api *api) OpenShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.OpenShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	_, err = api.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	if err != store.ErrNotFound {
		if err != nil {
			return nil, err
		}
		return nil, ErrShiftAlreadyExists
	}

	shift, err := rotation.makeShift(shiftNumber, nil)
	if err != nil {
		return nil, err
	}
	shift.ShiftStatus = store.ShiftStatusOpen

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	api.messageShiftOpened(rotation, shiftNumber, shift)
	logger.Infof("%s opened shift %s in %s.", MarkdownUser(api.actingUser), MarkdownShift(shiftNumber, shift), MarkdownRotation(rotation))
	return shift, nil
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
	err = api.expandShift(shift)
	if err != nil {
		return nil, err
	}
	return shift, nil
}
