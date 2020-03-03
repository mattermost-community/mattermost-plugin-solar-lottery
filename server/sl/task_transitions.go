// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/pkg/errors"
)

var ErrAlreadyExists = errors.New("already exists")

// func (sl *sl) StartTask(rotation *Rotation, t *Task) error {
// 	err := sl.Filter(
// 		withActingUserExpanded,
// 		withRotationExpanded(rotation),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	logger := sl.Logger.Timed().With(bot.LogContext{
// 		"Location":       "sl.StartTask",
// 		"ActingUsername": sl.actingUser.MattermostUsername(),
// 		"RotationID":     rotation.RotationID,
// 		"TaskID":         t.TaskID,
// 	})

// 	shift, err := sl.startShift(rotation, shiftNumber)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.Infof("%s started %s.", sl.actingUser.Markdown(), shift.Markdown())
// 	return shift, nil
// }

// func (sl *sl) DebugDeleteShift(rotation *Rotation, shiftNumber int) error {
// 	err := sl.Filter(
// 		withActingUserExpanded,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	logger := sl.Logger.Timed().With(bot.LogContext{
// 		"Location":       "sl.DebugDeleteShift",
// 		"ActingUsername": sl.actingUser.MattermostUsername(),
// 		"RotationID":     rotation.RotationID,
// 		"ShiftNumber":    shiftNumber,
// 	})

// 	err = sl.ShiftStore.DeleteShift(rotation.RotationID, shiftNumber)
// 	if err != nil {
// 		return err
// 	}

// 	logger.Infof("%s deleted shift %v in %s.", sl.actingUser.Markdown(), shiftNumber, rotation.Markdown())
// 	return nil
// }

// func (sl *sl) FinishShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
// 	err := sl.Filter(
// 		withActingUserExpanded,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logger := sl.Logger.Timed().With(bot.LogContext{
// 		"Location":       "sl.FinishShift",
// 		"ActingUsername": sl.actingUser.MattermostUsername(),
// 		"RotationID":     rotation.RotationID,
// 		"ShiftNumber":    shiftNumber,
// 	})

// 	shift, err := sl.finishShift(rotation, shiftNumber)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.Infof("%s finished %s.", sl.actingUser.Markdown(), shift.Markdown())
// 	return shift, nil
// }

// func (sl *sl) startTask(rotation *Rotation, t *Task) error {
// 	if t.Status == store.TaskStatusStarted {
// 		return errors.New("already started")
// 	}
// 	if t.Status != store.TaskStatusOpen {
// 		return errors.Errorf("can't start a shift which is %q, must be %q", t.Status, store.TaskStatusOpen)
// 	}

// 	for _, user := range rotation.ShiftUsers(shift) {
// 		rotation.markShiftUserServed(user, shiftNumber, shift)
// 		_, err = sl.storeUserWelcomeNew(user)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sl.messageShiftStarted(rotation, shift)
// 	return shift, nil
// }

func (sl *sl) finishTask(r *Rotation, t *Task) (*Task, error) {
	if t.Status == TaskStatusFinished {
		return t, nil
	}
	if t.Status != TaskStatusInProgress {
		return nil, errors.Errorf("can't finish a task which is %s, must be started", t.Status)
	}

	t.Status = TaskStatusFinished
	err := sl.Store.Entity(KeyTask).Store(t.TaskID, t)
	if err != nil {
		return nil, err
	}

	sl.messageTaskFinished(r, t)

	return t, nil
}

// func (sl *sl) transitionTaskToStatus(rotation *Rotation, t *Task, toStatus store.TaskStatus) error {
// 	fromStatus := t.Status
// 	fromTaskIDs, err := sl.RotationTasksStore.LoadRotationTaskIDs(rotation.RotationID, fromStatus)
// 	if err != nil {
// 		return err
// 	}
// 	toTaskIDs, err := sl.RotationTasksStore.LoadRotationTaskIDs(rotation.RotationID, toStatus)
// 	if err != nil {
// 		return err
// 	}
// 	for _, id := range toTaskIDs {
// 		if id == t.TaskID {
// 			return errors.Errorf("%s is already in %s", t.TaskID, toStatus)
// 		}
// 	}

// 	taskIDs := []string{}
// 	for _, id := range fromTaskIDs {
// 		if id != t.TaskID {
// 			taskIDs = append(taskIDs, id)
// 		}
// 	}
// 	err = sl.RotationTasksStore.StoreRotationTaskIDs(rotation.RotationID, fromStatus, taskIDs)
// 	if err != nil {
// 		return err
// 	}

// 	taskIDs = append(toTaskIDs, t.TaskID)
// 	err = sl.RotationTasksStore.StoreRotationTaskIDs(rotation.RotationID, toStatus, taskIDs)
// 	if err != nil {
// 		return err
// 	}

// 	t.Status = toStatus
// 	err = sl.TaskStore.StoreTask(t.Task)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
