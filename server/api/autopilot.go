// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
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
	finishedStatus := status(err)
	if err == nil {
		finishedStatus = status(err, fmt.Sprintf("finished %s.\n%s",
			rotation.ShiftRef(finishedShiftNumber),
			api.MarkdownIndent(finishedShift.MarkdownBullets(rotation), "  ")))
	}

	_, filledShifts, filledAdded, err :=
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
				"  - %s: **added users** %s.\n%s",
				shift.Markdown(),
				filledAdded[i].Markdown(),
				api.MarkdownIndent(shift.MarkdownBullets(rotation), "    "))
		} else {
			fillStatus += fmt.Sprintf("  - %s: no change.\n", shift.Markdown())
		}
	}

	_, startedShift, err := api.autopilotStartShift(rotation, now, currentShiftNumber)
	startedStatus := status(err)
	if err == nil {
		startedStatus = status(err, fmt.Sprintf("started %s.\n%s",
			startedShift.Markdown(),
			api.MarkdownIndent(startedShift.MarkdownBullets(rotation), "  ")))
	}

	currentNotified, err := api.autopilotNotifyCurrent(rotation, now, currentShiftNumber)
	currentNotifiedStatus := status(err, fmt.Sprintf("notified %s", currentNotified.Markdown()))

	nextNotified, err := api.autopilotNotifyNext(rotation, now, currentShiftNumber)
	nextNotifiedStatus := status(err, fmt.Sprintf("notified %s", nextNotified.Markdown()))

	logger.Infof("%s ran autopilot on %s for %v. Status:\n"+
		"- finish previous shift: %s\n"+
		"- fill shift(s): %s\n"+
		"- start next shift: %s\n"+
		"- notify current shift's users: %s\n"+
		"- notify next shift's users: %s\n",
		api.actingUser.Markdown(), rotation.Markdown(), now,
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
	finishShiftNumber, err := rotation.ShiftNumberForTime(now.Add(-1 * 24 * time.Hour))
	if err != nil {
		return 0, nil, err
	}
	if finishShiftNumber == -1 || finishShiftNumber == currentShiftNumber {
		return 0, nil, errors.New("no previous shift")
	}

	finishedShift, err := api.finishShift(rotation, finishShiftNumber)
	if err != nil {
		return 0, nil, err
	}
	return finishShiftNumber, finishedShift, nil
}

func (api *api) autopilotStartShift(rotation *Rotation, now time.Time, currentShiftNumber int) (int, *Shift, error) {
	if !rotation.Autopilot.StartFinish {
		return 0, nil, errors.New("not configured")
	}
	if currentShiftNumber == -1 {
		return 0, nil, errors.New("no shift to start")
	}

	currentShift, err := api.startShift(rotation, currentShiftNumber)
	if err != nil {
		return 0, nil, err
	}
	return currentShiftNumber, currentShift, nil
}

func (api *api) autopilotFill(rotation *Rotation, now time.Time, currentShiftNumber int, logger bot.Logger) ([]int, []*Shift, []UserMap, error) {
	if !rotation.Autopilot.Fill {
		return nil, nil, nil, errors.New("not configured to auto-fill")
	}

	startingShiftNumber := currentShiftNumber
	if startingShiftNumber < 0 {
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

	return api.fillShifts(rotation, startingShiftNumber, numShifts, now, logger)
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

	api.messageShiftWillFinish(rotation, currentShift)

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

	api.messageShiftWillStart(rotation, nextShift)

	nextShift.Shift.Autopilot.NotifiedStart = now
	err = api.ShiftStore.StoreShift(rotation.RotationID, nextShiftNumber, nextShift.Shift)
	if err != nil {
		return nil, err
	}

	return rotation.ShiftUsers(nextShift), nil
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
		if err != nil && err != ErrAlreadyExists {
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
			return filledShiftNumbers, filledShifts, addedUsers,
				errors.WithMessagef(err, "failed to join autofilled users to %s", loadedShift.Markdown())
		}

		err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, loadedShift.Shift)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers,
				errors.WithMessagef(err, "failed to store autofilled %s", loadedShift.Markdown())
		}

		api.messageShiftJoined(added, rotation, shift)
		appendShift(shiftNumber, shift, added)
	}

	return filledShiftNumbers, filledShifts, addedUsers, nil
}
