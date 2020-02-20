// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

// var ErrShiftMustBeOpen = errors.New("must be `open`")
// var ErrUserAlreadyInShift = errors.New("user is already in shift")

// func (sl *sl) AddUsersToTask(mattermostUsernames string, rotation *Rotation, t *Task) (UserMap, error) {
// 	err := sl.Filter(
// 		withActingUserExpanded,
// 		withMattermostUsersExpanded(mattermostUsernames),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logger := sl.Logger.Timed().With(bot.LogContext{
// 		"Location":            "sl.VolunteerUsers",
// 		"ActingUsername":      sl.actingUser.MattermostUsername(),
// 		"MattermostUsernames": mattermostUsernames,
// 		"RotationID":          rotation.RotationID,
// 		"TaskID":              t.TaskID,
// 	})

// 	added, err := sl.addUsersToTask(rotation, t, sl.users, true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = sl.TaskStore.StoreTask(t.Task)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sl.messageAddedUsersToTask(added, rotation, t)
// 	logger.Infof("%s volunteered %s to %s.",
// 		sl.actingUser.Markdown(), added.MarkdownWithSkills(), t.Markdown())
// 	return added, nil
// }

// func (sl *sl) IsTaskFilled(rotation *Rotation) (shift *Shift, ready bool, whyNot string, err error) {
// 	shift, err = sl.loadShift(rotation, shiftNumber)
// 	if err != nil {
// 		return nil, false, "", err
// 	}
// 	if shift.Status != store.ShiftStatusOpen {
// 		return nil, false, "", ErrShiftMustBeOpen
// 	}

// 	shiftUsers := rotation.ShiftUsers(shift)
// 	unmetNeeds := UnmetNeeds(rotation.Needs, shiftUsers)
// 	unmetCapacity := 0
// 	if rotation.Size != 0 {
// 		unmetCapacity = rotation.Size - len(shift.MattermostUserIDs)
// 	}

// 	if len(unmetNeeds) == 0 && unmetCapacity <= 0 {
// 		return shift, true, "", nil
// 	}

// 	whyNot = autofill.Error{
// 		UnmetNeeds:    unmetNeeds,
// 		UnmetCapacity: unmetCapacity,
// 		Err:           errors.New("not ready"),
// 		ShiftNumber:   shiftNumber,
// 	}.Error()

// 	return shift, false, whyNot, nil
// }

// func (sl *sl) FillTask(rotation *Rotation, t *Task) (UserMap, error) {
// 	err := sl.Filter(
// 		withActingUserExpanded,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger := sl.Logger.Timed().With(bot.LogContext{
// 		"Location":       "sl.FillShifts",
// 		"ActingUsername": sl.actingUser.MattermostUsername(),
// 		"RotationID":     rotation.RotationID,
// 		"TaskID":         t.TaskID,
// 	})

// 	_, shifts, addedUsers, err := sl.fillShifts(rotation, shiftNumber, 1, time.Time{}, logger)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	if len(shifts) == 0 || len(addedUsers) == 0 {
// 		logger.Infof("%s tried to fill %v, nothing to do.",
// 			sl.actingUser.Markdown(), rotation.ShiftRef(shiftNumber))
// 		return nil, nil, nil
// 	}

// 	shift := shifts[0]
// 	added := addedUsers[0]
// 	logger.Infof("%s filled %s, added %s.",
// 		sl.actingUser.Markdown(), shift.Markdown(), addedUsers[0].MarkdownWithSkills())
// 	return shift, added, nil
// }

// func (sl *sl) addUsersToTask(rotation *Rotation, t *Task, users UserMap, persist bool) (UserMap, error) {
// 	if t.Status != store.TaskStatusOpen {
// 		return nil, errors.Errorf("can't join task with status %s, must be %s", t.Status, store.TaskStatusOpen)
// 	}

// 	added := UserMap{}
// 	for _, user := range users {
// 		if t.MattermostUserIDs[user.MattermostUserID] != "" {
// 			continue
// 		}
// 		t.MattermostUserIDs.Add(user.MattermostUserID)
// 		added[user.MattermostUserID] = user
// 	}

// 	// err := sl.addEventToUsers(added, NewShiftEvent(rotation, shiftNumber, shift), persist)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	return added, nil
// }
