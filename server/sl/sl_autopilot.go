// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (sl *sl) autopilotRemindFinish(r *Rotation, now types.Time) (md.Markdowner, error) {
	if !r.AutopilotSettings.RemindFinish {
		return md.MD("task finish reminder: not configured"), nil
	}
	filtered := r.queryTasks(r.isAutopilotRemindFinish, now)
	if filtered.IsEmpty() {
		return md.MD("task finish reminder: nothing to do"), nil
	}

	var notified = NewUsers()
	for _, t := range filtered.AsArray() {
		for _, u := range t.Users.AsArray() {
			sl.dmUserTaskWillFinish(u, t)
			notified.Set(u)
		}

		t.AutopilotRemindedFinish = true
		err := sl.storeTask(t)
		if err != nil {
			return nil, err
		}
	}

	return md.Markdownf("task finish reminder: messaged %v users of %v tasks", notified.Len(), filtered.Len()), nil
}

func (sl *sl) autopilotRemindStart(r *Rotation, now types.Time) (md.Markdowner, error) {
	if !r.AutopilotSettings.RemindStart {
		return md.MD("task start reminder: not configured"), nil
	}
	filtered := r.queryTasks(r.isAutopilotRemindStart, now)
	if filtered.IsEmpty() {
		return md.MD("task start reminder: nothing to do"), nil
	}

	var notified = NewUsers()
	for _, t := range filtered.AsArray() {
		for _, u := range t.Users.AsArray() {
			sl.dmUserTaskWillStart(u, t)
			notified.Set(u)
		}

		t.AutopilotRemindedStart = true
		err := sl.storeTask(t)
		if err != nil {
			return nil, err
		}
	}

	return md.Markdownf("task finish reminder: messaged %v users of %v tasks", notified.Len(), filtered.Len()), nil
}

func (sl *sl) autopilotFinish(r *Rotation, now types.Time) (md.Markdowner, error) {
	if !r.AutopilotSettings.StartFinish {
		return md.MD("start/finish: not configured"), nil
	}
	filtered := r.queryTasks(r.isAutopilotFinish, now)
	if filtered.IsEmpty() {
		return md.MD("start/finish: nothing to do"), nil
	}

	for _, t := range filtered.AsArray() {
		err := sl.transitionTask(r, t, now, TaskStateFinished)
		if err != nil {
			return nil, err
		}
	}

	return md.Markdownf("finished: %v tasks", filtered.Len()), nil
}

func (sl *sl) autopilotStart(r *Rotation, now types.Time) (md.Markdowner, error) {
	if !r.AutopilotSettings.StartFinish {
		return md.MD("start/finish: not configured"), nil
	}
	filtered := r.queryTasks(r.isAutopilotStart, now)
	if filtered.IsEmpty() {
		return md.MD("start/finish: nothing to do"), nil
	}

	for _, t := range filtered.AsArray() {
		err := sl.transitionTask(r, t, now, TaskStateStarted)
		if err != nil {
			return nil, err
		}
	}

	return md.Markdownf("started: %v tasks", filtered.Len()), nil
}

func (s *sl) autopilotFillSchedule(r *Rotation, now types.Time) (md.Markdowner, error) {
	if !r.AutopilotSettings.Schedule {
		return md.MD("fill and schedule: not configured"), nil
	}
	filtered := r.queryTasks(r.isAutopilotSchedule, now)
	if filtered.IsEmpty() {
		return md.MD("fill and schedule: nothing to do"), nil
	}

	var messages []md.Markdowner
	for _, t := range filtered.AsArray() {
		outFill, err := s.FillTask(InAssignTask{
			TaskID: t.TaskID,
			Time:   now,
		})
		if err != nil {
			return nil, err
		}

		outTransition, err := s.TransitionTask(InTransitionTask{
			TaskID: t.TaskID,
			State:  TaskStateScheduled,
			Time:   now,
		})
		if err != nil {
			return nil, err
		}
		messages = append(messages, outFill.Markdown(), md.MD(", "), outTransition.Markdown(), md.MD("\n"))
	}

	text := md.Markdownf("fill and schedule: processed %v shifts:\n", len(messages))
	for _, m := range messages {
		text += m.Markdown()
	}
	return text, nil
}

func (sl *sl) autopilotCreate(r *Rotation, now types.Time) (md.Markdowner, error) {
	if r.TaskType == TaskTypeTicket {
		return md.MD("create shift: tickets can not be auto-created"), nil
	}
	if !r.AutopilotSettings.Create {
		return md.MD("create shift: not configured"), nil
	}

	var messages []md.Markdowner
	period := r.TaskSettings.ShiftPeriod
	upTo := now.Add(r.AutopilotSettings.CreatePrior)
	for num, start := period.ForTime(r.Beginning, now); start.Before(upTo); num, start = num+1, period.ForNumber(start, 1) {
		exists := r.queryTasks(r.isPendingForTime, start)
		if !exists.IsEmpty() {
			continue
		}
		out, err := sl.CreateShift(InCreateShift{
			RotationID: r.RotationID,
			Number:     num,
		})
		if err != nil {
			return nil, err
		}
		messages = append(messages, out, md.MD("\n"))
	}

	if len(messages) == 0 {
		return md.MD("create shift: nothing to do"), nil
	}
	text := md.Markdownf("create shift: created %v shifts:\n", len(messages))
	for _, m := range messages {
		text += m.Markdown()
	}
	return text, nil
}

// func (sl *sl) autopilotFill(rotation *Rotation, now time.Time, currentShiftNumber int, logger bot.Logger) ([]int, []*Shift, []Users, error) {
// 	if !rotation.Autopilot.Fill {
// 		return nil, nil, nil, errors.New("not configured to auto-fill")
// 	}

// 	startingShiftNumber := currentShiftNumber
// 	if startingShiftNumber < 0 {
// 		startingShiftNumber = 0
// 	}

// 	upToShiftNumber, err := rotation.ShiftNumberForTime(now.Add(rotation.Autopilot.FillPrior))
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	if upToShiftNumber < startingShiftNumber {
// 		return nil, nil, nil, errors.New("nothing to do")
// 	}
// 	numShifts := upToShiftNumber - startingShiftNumber + 1

// 	return sl.fillShifts(rotation, startingShiftNumber, numShifts, now, logger)
// }

// func (sl *sl) fillShifts(rotation *Rotation, startingShiftNumber, numShifts int, now time.Time, logger bot.Logger) ([]int, []*Shift, []Users, error) {
// 	// Guess' logs are too verbose - suppress
// 	prevLogger := sl.Logger
// 	sl.Logger = &bot.NilLogger{}
// 	shifts, err := sl.Guess(rotation, startingShiftNumber, numShifts)
// 	sl.Logger = prevLogger
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	if len(shifts) != numShifts {
// 		return nil, nil, nil, errors.New("unreachable, must match")
// 	}

// 	var filledShiftNumbers []int
// 	var filledShifts []*Shift
// 	var addedUsers []Users

// 	appendShift := func(shiftNumber int, shift *Shift, added Users) {
// 		filledShiftNumbers = append(filledShiftNumbers, shiftNumber)
// 		filledShifts = append(filledShifts, shift)
// 		addedUsers = append(addedUsers, added)
// 	}

// 	shiftNumber := startingShiftNumber - 1
// 	for n := 0; n < numShifts; n++ {
// 		shiftNumber++

// 		loadedShift, err := sl.OpenShift(rotation, shiftNumber)
// 		if err != nil && err != ErrAlreadyExists {
// 			return nil, nil, nil, err
// 		}
// 		if loadedShift.Status != store.ShiftStatusOpen {
// 			appendShift(shiftNumber, loadedShift, nil)
// 			continue
// 		}
// 		if !loadedShift.Autopilot.Filled.IsZero() {
// 			appendShift(shiftNumber, loadedShift, nil)
// 			continue
// 		}

// 		before := rotation.ShiftUsers(loadedShift).Clone(false)

// 		// shifts coming from Guess are either loaded with their respective
// 		// status, or are Open. (in reality should always be Open).
// 		shift := shifts[n]
// 		added := NewUsers()
// 		for id, user := range rotation.ShiftUsers(shift) {
// 			if before[id] == nil {
// 				added[id] = user
// 			}
// 		}
// 		if len(added) == 0 {
// 			appendShift(shiftNumber, loadedShift, nil)
// 			continue
// 		}

// 		loadedShift.Autopilot.Filled = now

// 		_, err = sl.joinShift(rotation, shiftNumber, loadedShift, added, true)
// 		if err != nil {
// 			return filledShiftNumbers, filledShifts, addedUsers,
// 				errors.WithMessagef(err, "failed to join autofilled users to %s", loadedShift.Markdown())
// 		}

// 		err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, loadedShift.Shift)
// 		if err != nil {
// 			return filledShiftNumbers, filledShifts, addedUsers,
// 				errors.WithMessagef(err, "failed to store autofilled %s", loadedShift.Markdown())
// 		}

// 		sl.messageShiftJoined(added, rotation, shift)
// 		appendShift(shiftNumber, shift, added)
// 	}

// 	return filledShiftNumbers, filledShifts, addedUsers, nil
// }
