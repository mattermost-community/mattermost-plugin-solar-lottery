// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) AutopilotRotation(rotation *Rotation, now time.Time) error {
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.AutopilotRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"Time":           now,
	})

	if !rotation.Autopilot.On {
		return nil
	}
	currentShiftNumber, err := rotation.ShiftNumberForTime(now)
	if err != nil {
		return err
	}

	api.Debugf("<><> currentShiftNumber: %v", currentShiftNumber)
	err = api.autopilotFinish(rotation, now, currentShiftNumber)
	if err != nil {
		api.Logger.Infof("Failed to finish previous shift: %s", err.Error())
	}

	err = api.autopilotStart(rotation, now, currentShiftNumber)
	if err != nil {
		api.Logger.Infof("Failed to start next shift: %s", err.Error())
	}

	err = api.autopilotFill(rotation, now, currentShiftNumber, logger)
	if err != nil {
		api.Logger.Infof("Failed to fill shift(s): %s", err.Error())
	}

	err = api.autopilotNotifyCurrent(rotation, now, currentShiftNumber)
	if err != nil {
		api.Logger.Infof("Failed to notify current users: %s", err.Error())
	}

	err = api.autopilotNotifyNext(rotation, now, currentShiftNumber)
	if err != nil {
		api.Logger.Infof("Failed to notify next shift's users: %s", err.Error())
	}

	logger.Infof("%s ran autopilot on %s.", api.MarkdownUser(api.actingUser), MarkdownRotation(rotation))
	return nil
}

func (api *api) autopilotFinish(rotation *Rotation, now time.Time, currentShiftNumber int) error {
	if !rotation.Autopilot.StartFinish {
		return nil
	}
	prevShiftNumber, err := rotation.ShiftNumberForTime(now.Add(-1 * 24 * time.Hour))
	if err != nil {
		return err
	}
	if prevShiftNumber == -1 || prevShiftNumber == currentShiftNumber {
		api.Debugf("<><> Finish: nothing to do")
		return nil
	}

	_, err = api.finishShift(rotation, prevShiftNumber)
	if err != nil {
		return err
	}

	return err
}

func (api *api) autopilotStart(rotation *Rotation, now time.Time, currentShiftNumber int) error {
	if !rotation.Autopilot.StartFinish || currentShiftNumber == -1 {
		api.Debugf("<><> Start: nothing to do 1")
		return nil
	}

	prevShiftNumber, err := rotation.ShiftNumberForTime(now.Add(-1 * 24 * time.Hour))
	if err != nil {
		return err
	}
	if prevShiftNumber == currentShiftNumber {
		api.Debugf("<><> Start: nothing to do 2")
		return nil
	}

	_, err = api.startShift(rotation, currentShiftNumber)
	if err != nil {
		return err
	}

	return err
}

func (api *api) autopilotFill(rotation *Rotation, now time.Time, currentShiftNumber int, logger bot.Logger) error {
	if !rotation.Autopilot.Fill {
		return nil
	}

	startingShiftNumber := currentShiftNumber
	if currentShiftNumber < 0 {
		startingShiftNumber = 0
	}

	upToShiftNumber, err := rotation.ShiftNumberForTime(now.Add(rotation.Autopilot.FillPrior))
	if err != nil {
		return err
	}
	if upToShiftNumber < currentShiftNumber {
		return nil
	}
	numShifts := upToShiftNumber - currentShiftNumber + 1

	_, _, err = api.fillShifts(rotation, startingShiftNumber, numShifts, true, now, logger)
	if err != nil {
		return err
	}

	return nil
}

func (api *api) autopilotNotifyCurrent(rotation *Rotation, now time.Time, currentShiftNumber int) error {
	if !rotation.Autopilot.Notify || currentShiftNumber == -1 {
		return nil
	}
	_, e, err := rotation.ShiftDatesForNumber(currentShiftNumber)
	if err != nil {
		return nil
	}
	if e.After(now.Add(rotation.Autopilot.NotifyPrior)) {
		return nil
	}

	currentShift, err := api.loadShift(rotation, currentShiftNumber)
	if err != nil {
		return err
	}
	api.Debugf("<><> NotifyCurrent:  %v", currentShiftNumber)
	if !currentShift.Autopilot.NotifiedFinish.IsZero() {
		return nil
	}

	api.messageShiftWillFinish(rotation, currentShiftNumber, currentShift)

	currentShift.Shift.Autopilot.NotifiedFinish = now
	return api.ShiftStore.StoreShift(rotation.RotationID, currentShiftNumber, currentShift.Shift)
}

func (api *api) autopilotNotifyNext(rotation *Rotation, now time.Time, currentShiftNumber int) error {
	if !rotation.Autopilot.Notify {
		return nil
	}
	nextShiftNumber := currentShiftNumber + 1
	s, _, err := rotation.ShiftDatesForNumber(nextShiftNumber)
	if err != nil {
		return nil
	}
	if s.After(now.Add(rotation.Autopilot.NotifyPrior)) {
		return nil
	}

	nextShift, err := api.loadShift(rotation, nextShiftNumber)
	if err != nil {
		return err
	}
	api.Debugf("<><> NotifyNext:  %v", nextShiftNumber)
	if !nextShift.Autopilot.NotifiedStart.IsZero() {
		return nil
	}

	api.messageShiftWillStart(rotation, nextShiftNumber, nextShift)

	nextShift.Shift.Autopilot.NotifiedStart = now
	return api.ShiftStore.StoreShift(rotation.RotationID, nextShiftNumber, nextShift.Shift)
}
