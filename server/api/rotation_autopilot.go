// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) AutopilotRotation(rotation *Rotation, now time.Time) error {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return err
	}
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

	api.Debugf("running autopilot for shiftNumber: %v", currentShiftNumber)
	status := func(err error, message ...string) string {
		switch {
		case err == nil && len(message) > 0 && message[0] != "":
			return "ok: " + message[0]
		case err == nil && (len(message) == 0 || message[0] == ""):
			return "ok"
		case err != nil:
			return "**" + err.Error() + "**"
		}
		return "unknown"
	}

	finishedShiftNumber, finishedShift, err :=
		api.autopilotFinishShift(rotation, now, currentShiftNumber)
	finishedStatus := status(err, fmt.Sprintf("finished %s",
		api.MarkdownShiftWithDetails(rotation, finishedShiftNumber, finishedShift)))

	filledShiftNumbers, filledShifts, filledAdded, err :=
		api.autopilotFill(rotation, now, currentShiftNumber, logger)
	fillStatus := fmt.Sprintf("ok: **processed %v shifts**", len(filledShifts))
	if err != nil {
		if len(filledShifts) > 0 {
			fillStatus = fmt.Sprintf("error: processed %v shifts, then **failed: %v**.", len(filledShifts), err)
		} else {
			fillStatus = err.Error()
		}
	}
	if len(filledShifts) > 0 {
		fillStatus += "\n"
	}
	for i, shift := range filledShifts {
		if len(filledAdded[i]) > 0 {
			fillStatus += fmt.Sprintf(
				"  - %s: **added users** %s.\n"+
					"    - %s\n",
				MarkdownShift(rotation, filledShiftNumbers[i]),
				api.MarkdownUsers(filledAdded[i]),
				api.MarkdownShiftDetails(rotation, filledShiftNumbers[i], shift))
		} else {
			fillStatus += fmt.Sprintf(
				"  - %s: no change.\n",
				MarkdownShift(rotation, filledShiftNumbers[i]))
		}
	}

	startedShiftNumber, startedShift, err := api.autopilotStartShift(rotation, now, currentShiftNumber)
	startedStatus := status(err, fmt.Sprintf("started %s",
		api.MarkdownShiftWithDetails(rotation, startedShiftNumber, startedShift)))

	currentNotified, err := api.autopilotNotifyCurrent(rotation, now, currentShiftNumber)
	currentNotifiedStatus := status(err, fmt.Sprintf("notified %s", api.MarkdownUsers(currentNotified)))

	nextNotified, err := api.autopilotNotifyNext(rotation, now, currentShiftNumber)
	nextNotifiedStatus := status(err, fmt.Sprintf("notified %s", api.MarkdownUsers(nextNotified)))

	logger.Infof("%s ran autopilot on %s for %v. Status:\n"+
		"- finish previous shift: %s\n"+
		"- fill shift(s): %s\n"+
		"- start next shift: %s\n"+
		"- notify current shift's users: %s\n"+
		"- notify next shift's users: %s\n",
		api.MarkdownUser(api.actingUser), MarkdownRotation(rotation), now,
		finishedStatus,
		fillStatus,
		startedStatus,
		currentNotifiedStatus,
		nextNotifiedStatus,
	)

	return nil
}

func (api *api) autopilotFinishShift(rotation *Rotation, now time.Time, currentShiftNumber int) (int, *Shift, error) {
	if !rotation.Autopilot.StartFinish {
		return 0, nil, errors.New("not configured")
	}
	prevShiftNumber, err := rotation.ShiftNumberForTime(now.Add(-1 * 24 * time.Hour))
	if err != nil {
		return 0, nil, err
	}
	if prevShiftNumber == -1 || prevShiftNumber == currentShiftNumber {
		return 0, nil, errors.New("no previous shift")
	}

	shift, err := api.finishShift(rotation, prevShiftNumber)
	if err != nil {
		return 0, nil, err
	}
	return prevShiftNumber, shift, nil
}

func (api *api) autopilotStartShift(rotation *Rotation, now time.Time, currentShiftNumber int) (int, *Shift, error) {
	if !rotation.Autopilot.StartFinish {
		return 0, nil, errors.New("not configured")
	}
	if currentShiftNumber == -1 {
		return 0, nil, errors.New("no shift to start")
	}

	shift, err := api.startShift(rotation, currentShiftNumber)
	if err != nil {
		return 0, nil, err
	}
	return currentShiftNumber, shift, nil
}

func (api *api) autopilotFill(rotation *Rotation, now time.Time, currentShiftNumber int, logger bot.Logger) ([]int, []*Shift, []UserMap, error) {
	if !rotation.Autopilot.Fill {
		return nil, nil, nil, errors.New("not configured to auto-fill")
	}

	startingShiftNumber := currentShiftNumber
	if currentShiftNumber < 0 {
		startingShiftNumber = 0
	}

	upToShiftNumber, err := rotation.ShiftNumberForTime(now.Add(rotation.Autopilot.FillPrior))
	if err != nil {
		return nil, nil, nil, err
	}
	if upToShiftNumber < startingShiftNumber {
		return nil, nil, nil, errors.New("nothing to do")
	}
	numShifts := upToShiftNumber - startingShiftNumber + 1

	return api.fillShifts(rotation, startingShiftNumber, numShifts, true, now, logger)
}

func (api *api) autopilotNotifyCurrent(rotation *Rotation, now time.Time, currentShiftNumber int) (UserMap, error) {
	if !rotation.Autopilot.Notify {
		return nil, errors.New("not configured")
	}
	if currentShiftNumber == -1 {
		return nil, errors.New("no shift")
	}
	_, e, err := rotation.ShiftDatesForNumber(currentShiftNumber)
	if err != nil {
		return nil, err
	}
	if e.After(now.Add(rotation.Autopilot.NotifyPrior)) {
		return nil, errors.New("nothing to do")
	}

	currentShift, err := api.loadShift(rotation, currentShiftNumber)
	if err != nil {
		return nil, err
	}
	if !currentShift.Autopilot.NotifiedFinish.IsZero() {
		return nil, errors.New("already notified")
	}

	api.messageShiftWillFinish(rotation, currentShiftNumber, currentShift)

	currentShift.Shift.Autopilot.NotifiedFinish = now
	err = api.ShiftStore.StoreShift(rotation.RotationID, currentShiftNumber, currentShift.Shift)
	if err != nil {
		return nil, err
	}

	return rotation.ShiftUsers(currentShift), nil
}

func (api *api) autopilotNotifyNext(rotation *Rotation, now time.Time, currentShiftNumber int) (UserMap, error) {
	if !rotation.Autopilot.Notify {
		return nil, errors.New("not configured to notify")
	}
	nextShiftNumber := currentShiftNumber + 1
	s, _, err := rotation.ShiftDatesForNumber(nextShiftNumber)
	if err != nil {
		return nil, err
	}
	if s.After(now.Add(rotation.Autopilot.NotifyPrior)) {
		return nil, errors.New("nothing to do")
	}

	nextShift, err := api.loadShift(rotation, nextShiftNumber)
	if err != nil {
		return nil, err
	}
	if !nextShift.Autopilot.NotifiedStart.IsZero() {
		return nil, errors.New("already notified")
	}

	api.messageShiftWillStart(rotation, nextShiftNumber, nextShift)

	nextShift.Shift.Autopilot.NotifiedStart = now
	err = api.ShiftStore.StoreShift(rotation.RotationID, nextShiftNumber, nextShift.Shift)
	if err != nil {
		return nil, err
	}

	return rotation.ShiftUsers(nextShift), nil
}
