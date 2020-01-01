// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"
)

func (api *api) AutopilotRotation(rotation *Rotation, now time.Time) error {
	if !rotation.Autopilot.On {
		return nil
	}
	currentShiftNumber, err := rotation.ShiftNumberForTime(now)
	if err != nil {
		return err
	}

	err = api.autopilotFinish(rotation, now, currentShiftNumber)
	if err != nil {
		api.Logger.Infof("Failed to finish previous shift: %s", err.Error())
	}

	err = api.autopilotStart(rotation, now, currentShiftNumber)
	if err != nil {
		api.Logger.Infof("Failed to start next shift: %s", err.Error())
	}

	err = api.autopilotFill(rotation, now, currentShiftNumber)
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
		return nil
	}

	prevShift, err := api.loadShift(rotation, prevShiftNumber)
	if err != nil {
		return err
	}
	if !prevShift.Autopilot.Finished.IsZero() {
		return nil
	}

	err = api.finishShift(rotation, prevShiftNumber, prevShift)
	if err != nil {
		return err
	}

	prevShift.Autopilot.Finished = now
	return api.ShiftStore.StoreShift(rotation.RotationID, prevShiftNumber, prevShift.Shift)
}

func (api *api) autopilotStart(rotation *Rotation, now time.Time, currentShiftNumber int) error {
	if !rotation.Autopilot.StartFinish || currentShiftNumber == -1 {
		return nil
	}

	prevShiftNumber, err := rotation.ShiftNumberForTime(now.Add(-1 * 24 * time.Hour))
	if err != nil {
		return err
	}
	if prevShiftNumber == currentShiftNumber {
		return nil
	}

	currentShift, err := api.loadShift(rotation, currentShiftNumber)
	if err != nil {
		return err
	}
	if !currentShift.Autopilot.Started.IsZero() {
		return nil
	}

	err = api.startShift(rotation, currentShiftNumber, currentShift)
	if err != nil {
		return err
	}

	currentShift.Autopilot.Started = now
	return api.ShiftStore.StoreShift(rotation.RotationID, currentShiftNumber, currentShift.Shift)
}

func (api *api) autopilotFill(rotation *Rotation, now time.Time, currentShiftNumber int) error {
	if !rotation.Autopilot.Fill {
		return nil
	}

	fillShiftNumber, err := rotation.ShiftNumberForTime(now.Add(rotation.Autopilot.FillPrior))
	if err != nil {
		return err
	}
	if fillShiftNumber == -1 {
		return nil
	}

	for n := currentShiftNumber; n <= fillShiftNumber; n++ {
		var shift *Shift
		shift, err = api.OpenShift(rotation, fillShiftNumber)
		if err != nil && err != ErrShiftAlreadyExists {
			return err
		}

		if !shift.Autopilot.Filled.IsZero() {
			continue
		}

		shift, _, _, _, _, err = api.FillShift(rotation, fillShiftNumber, true)
		if err != nil && err != ErrShiftMustBeOpen {
			return err
		}

		shift.Shift.Autopilot.Filled = now
		err = api.ShiftStore.StoreShift(rotation.RotationID, n, shift.Shift)
		if err != nil {
			return errors.WithMessagef(err, "failed to store autofilled %s", MarkdownShift(rotation, n, shift))
		}
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
	if !nextShift.Autopilot.NotifiedStart.IsZero() {
		return nil
	}

	api.messageShiftWillStart(rotation, nextShiftNumber, nextShift)

	nextShift.Shift.Autopilot.NotifiedStart = now
	return api.ShiftStore.StoreShift(rotation.RotationID, nextShiftNumber, nextShift.Shift)
}
