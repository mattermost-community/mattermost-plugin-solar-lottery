// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

var ErrAlreadyExists = errors.New("already exists")

func (sl *solarLottery) OpenShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.OpenShift",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := sl.loadShift(rotation, shiftNumber)
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

	err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	sl.messageShiftOpened(rotation, shift)
	logger.Infof("%s opened %s.", sl.actingUser.Markdown(), shift.Markdown())
	return shift, nil
}

func (sl *solarLottery) StartShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := sl.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.StartShift",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := sl.startShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}

	logger.Infof("%s started %s.", sl.actingUser.Markdown(), shift.Markdown())
	return shift, nil
}

func (sl *solarLottery) DebugDeleteShift(rotation *Rotation, shiftNumber int) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.DebugDeleteShift",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	err = sl.ShiftStore.DeleteShift(rotation.RotationID, shiftNumber)
	if err != nil {
		return err
	}

	logger.Infof("%s deleted shift %v in %s.", sl.actingUser.Markdown(), shiftNumber, rotation.Markdown())
	return nil
}

func (sl *solarLottery) FinishShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.FinishShift",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := sl.finishShift(rotation, shiftNumber)
	if err != nil {
		return nil, err
	}

	logger.Infof("%s finished %s.", sl.actingUser.Markdown(), shift.Markdown())
	return shift, nil
}

func (sl *solarLottery) startShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := sl.loadShift(rotation, shiftNumber)
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
		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return nil, err
		}
	}

	err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	sl.messageShiftStarted(rotation, shift)
	return shift, nil
}

func (sl *solarLottery) finishShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := sl.loadShift(rotation, shiftNumber)
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
	err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	sl.messageShiftFinished(rotation, shift)

	return shift, nil
}
