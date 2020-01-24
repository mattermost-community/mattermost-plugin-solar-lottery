// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Autopilot interface {
	AutopilotRotation(rotation *Rotation, now time.Time) error
}

func (sl *solarLottery) AutopilotRotation(rotation *Rotation, now time.Time) error {
	err := sl.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.AutopilotRotation",
		"ActingUsername": sl.actingUser.MattermostUsername(),
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

	sl.Debugf("running autopilot for shiftNumber: %v", currentShiftNumber)
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
		sl.autopilotFinishShift(rotation, now, currentShiftNumber)
	finishedStatus := status(err)
	if err == nil {
		finishedStatus = status(err, fmt.Sprintf("finished %s.\n%s",
			rotation.ShiftRef(finishedShiftNumber),
			utils.Indent(finishedShift.MarkdownBullets(rotation), "  ")))
	}

	_, filledShifts, filledAdded, err :=
		sl.autopilotFill(rotation, now, currentShiftNumber, logger)
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
				utils.Indent(shift.MarkdownBullets(rotation), "    "))
		} else {
			fillStatus += fmt.Sprintf("  - %s: no change.\n", shift.Markdown())
		}
	}

	_, startedShift, err := sl.autopilotStartShift(rotation, now, currentShiftNumber)
	startedStatus := status(err)
	if err == nil {
		startedStatus = status(err, fmt.Sprintf("started %s.\n%s",
			startedShift.Markdown(),
			utils.Indent(startedShift.MarkdownBullets(rotation), "  ")))
	}

	currentNotified, err := sl.autopilotNotifyCurrent(rotation, now, currentShiftNumber)
	currentNotifiedStatus := status(err, fmt.Sprintf("notified %s", currentNotified.Markdown()))

	nextNotified, err := sl.autopilotNotifyNext(rotation, now, currentShiftNumber)
	nextNotifiedStatus := status(err, fmt.Sprintf("notified %s", nextNotified.Markdown()))

	logger.Infof("%s ran autopilot on %s for %v. Status:\n"+
		"- finish previous shift: %s\n"+
		"- fill shift(s): %s\n"+
		"- start next shift: %s\n"+
		"- notify current shift's users: %s\n"+
		"- notify next shift's users: %s\n",
		sl.actingUser.Markdown(), rotation.Markdown(), now,
		finishedStatus,
		fillStatus,
		startedStatus,
		currentNotifiedStatus,
		nextNotifiedStatus,
	)

	return nil
}

func (sl *solarLottery) autopilotFinishShift(rotation *Rotation, now time.Time, currentShiftNumber int) (int, *Shift, error) {
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

	finishedShift, err := sl.finishShift(rotation, finishShiftNumber)
	if err != nil {
		return 0, nil, err
	}
	return finishShiftNumber, finishedShift, nil
}

func (sl *solarLottery) autopilotStartShift(rotation *Rotation, now time.Time, currentShiftNumber int) (int, *Shift, error) {
	if !rotation.Autopilot.StartFinish {
		return 0, nil, errors.New("not configured")
	}
	if currentShiftNumber == -1 {
		return 0, nil, errors.New("no shift to start")
	}

	currentShift, err := sl.startShift(rotation, currentShiftNumber)
	if err != nil {
		return 0, nil, err
	}
	return currentShiftNumber, currentShift, nil
}

func (sl *solarLottery) autopilotFill(rotation *Rotation, now time.Time, currentShiftNumber int, logger bot.Logger) ([]int, []*Shift, []UserMap, error) {
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

	return sl.fillShifts(rotation, startingShiftNumber, numShifts, now, logger)
}

func (sl *solarLottery) autopilotNotifyCurrent(rotation *Rotation, now time.Time, currentShiftNumber int) (UserMap, error) {
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

	currentShift, err := sl.loadShift(rotation, currentShiftNumber)
	if err != nil {
		return nil, err
	}
	if !currentShift.Autopilot.NotifiedFinish.IsZero() {
		return nil, errors.New("already notified")
	}

	sl.messageShiftWillFinish(rotation, currentShift)

	currentShift.Shift.Autopilot.NotifiedFinish = now
	err = sl.ShiftStore.StoreShift(rotation.RotationID, currentShiftNumber, currentShift.Shift)
	if err != nil {
		return nil, err
	}

	return rotation.ShiftUsers(currentShift), nil
}

func (sl *solarLottery) autopilotNotifyNext(rotation *Rotation, now time.Time, currentShiftNumber int) (UserMap, error) {
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

	nextShift, err := sl.loadShift(rotation, nextShiftNumber)
	if err != nil {
		return nil, err
	}
	if !nextShift.Autopilot.NotifiedStart.IsZero() {
		return nil, errors.New("already notified")
	}

	sl.messageShiftWillStart(rotation, nextShift)

	nextShift.Shift.Autopilot.NotifiedStart = now
	err = sl.ShiftStore.StoreShift(rotation.RotationID, nextShiftNumber, nextShift.Shift)
	if err != nil {
		return nil, err
	}

	return rotation.ShiftUsers(nextShift), nil
}

func (sl *solarLottery) fillShifts(rotation *Rotation, startingShiftNumber, numShifts int, now time.Time, logger bot.Logger) ([]int, []*Shift, []UserMap, error) {
	// Guess' logs are too verbose - suppress
	prevLogger := sl.Logger
	sl.Logger = &bot.NilLogger{}
	shifts, err := sl.Guess(rotation, startingShiftNumber, numShifts)
	sl.Logger = prevLogger
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

		loadedShift, err := sl.OpenShift(rotation, shiftNumber)
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

		_, err = sl.joinShift(rotation, shiftNumber, loadedShift, added, true)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers,
				errors.WithMessagef(err, "failed to join autofilled users to %s", loadedShift.Markdown())
		}

		err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, loadedShift.Shift)
		if err != nil {
			return filledShiftNumbers, filledShifts, addedUsers,
				errors.WithMessagef(err, "failed to store autofilled %s", loadedShift.Markdown())
		}

		sl.messageShiftJoined(added, rotation, shift)
		appendShift(shiftNumber, shift, added)
	}

	return filledShiftNumbers, filledShifts, addedUsers, nil
}
