// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

var ErrAlreadyExists = errors.New("already exists")

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

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != store.ErrNotFound {
		if err != nil {
			return nil, err
		}
		return shift, ErrAlreadyExists
	}

	shift, err = rotation.makeShift(shiftNumber)
	if err != nil {
		return nil, err
	}
	shift.Status = store.ShiftStatusOpen

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	api.messageShiftOpened(rotation, shift)
	logger.Infof("%s opened %s.", api.actingUser.Markdown(), shift.Markdown())
	return shift, nil
}

func (api *api) StartShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.StartShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := api.startShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}

	logger.Infof("%s started %s.", api.actingUser.Markdown(), shift.Markdown())
	return shift, nil
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

	logger.Infof("%s deleted shift %v in %s.", api.actingUser.Markdown(), shiftNumber, rotation.Markdown())
	return nil
}

func (api *api) FinishShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.FinishShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := api.finishShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}

	logger.Infof("%s finished %s.", api.actingUser.Markdown(), shift.Markdown())
	return shift, nil
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

	api.messageShiftStarted(rotation, shift)
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

	api.messageShiftFinished(rotation, shift)

	return shift, nil
}
