// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

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

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}

	err = api.finishShift(rotation, shiftNumber, shift)
	if err != nil {
		return nil, err
	}

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	logger.Infof("%s finished %s.", MarkdownUser(api.actingUser), MarkdownShift(rotation, shiftNumber, shift))
	return shift, nil
}

func (api *api) finishShift(rotation *Rotation, shiftNumber int, shift *Shift) error {
	if shift.Status == store.ShiftStatusFinished {
		return nil
	}
	if shift.Status != store.ShiftStatusStarted {
		return errors.Errorf("can't finish a shift which is %s, must be started", shift.Status)
	}

	shift.Status = store.ShiftStatusFinished
	api.messageShiftFinished(rotation, shiftNumber, shift)
	return nil
}
