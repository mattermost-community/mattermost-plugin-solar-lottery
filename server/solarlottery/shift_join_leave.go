// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/autofill"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

var ErrShiftMustBeOpen = errors.New("must be `open`")
var ErrUserAlreadyInShift = errors.New("user is already in shift")

func (sl *solarLottery) JoinShift(mattermostUsernames string, rotation *Rotation, shiftNumber int) (*Shift, UserMap, error) {
	err := sl.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return nil, nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":            "sl.VolunteerUsers",
		"ActingUsername":      sl.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"RotationID":          rotation.RotationID,
		"ShiftNumber":         shiftNumber,
	})

	shift, err := sl.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, nil, errors.Errorf("failed to load shift %v for rotation %s", shiftNumber, rotation.RotationID)
	}
	joined, err := sl.joinShift(rotation, shiftNumber, shift, sl.users, true)
	if err != nil {
		return nil, nil, err
	}
	err = sl.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, nil, err
	}

	sl.messageShiftJoined(joined, rotation, shift)
	logger.Infof("%s volunteered %s to %s.",
		sl.actingUser.Markdown(), joined.MarkdownWithSkills(), shift.Markdown())
	return shift, joined, nil
}
func (sl *solarLottery) IsShiftReady(rotation *Rotation, shiftNumber int) (shift *Shift, ready bool, whyNot string, err error) {
	shift, err = sl.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, false, "", err
	}
	if shift.Status != store.ShiftStatusOpen {
		return nil, false, "", ErrShiftMustBeOpen
	}

	shiftUsers := rotation.ShiftUsers(shift)
	unmetNeeds := UnmetNeeds(rotation.Needs, shiftUsers)
	unmetCapacity := 0
	if rotation.Size != 0 {
		unmetCapacity = rotation.Size - len(shift.MattermostUserIDs)
	}

	if len(unmetNeeds) == 0 && unmetCapacity <= 0 {
		return shift, true, "", nil
	}

	whyNot = autofill.Error{
		UnmetNeeds:    unmetNeeds,
		UnmetCapacity: unmetCapacity,
		Err:           errors.New("not ready"),
		ShiftNumber:   shiftNumber,
	}.Error()

	return shift, false, whyNot, nil
}

func (sl *solarLottery) FillShift(rotation *Rotation, shiftNumber int) (*Shift, UserMap, error) {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, nil, err
	}

	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.FillShifts",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	_, shifts, addedUsers, err := sl.fillShifts(rotation, shiftNumber, 1, time.Time{}, logger)
	if err != nil {
		return nil, nil, err
	}

	if len(shifts) == 0 || len(addedUsers) == 0 {
		logger.Infof("%s tried to fill %v, nothing to do.",
			sl.actingUser.Markdown(), rotation.ShiftRef(shiftNumber))
		return nil, nil, nil
	}

	shift := shifts[0]
	added := addedUsers[0]
	logger.Infof("%s filled %s, added %s.",
		sl.actingUser.Markdown(), shift.Markdown(), addedUsers[0].MarkdownWithSkills())
	return shift, added, nil
}

func (sl *solarLottery) joinShift(rotation *Rotation, shiftNumber int, shift *Shift, users UserMap, persist bool) (UserMap, error) {
	if shift.Status != store.ShiftStatusOpen {
		return nil, errors.Errorf("can't join a shift with status %s, must be Open", shift.Status)
	}

	joined := UserMap{}
	for _, user := range users {
		if shift.Shift.MattermostUserIDs[user.MattermostUserID] != "" {
			continue
		}
		if len(shift.MattermostUserIDs) >= rotation.Size {
			return nil, errors.Errorf("rotation size %v exceeded", rotation.Size)
		}
		shift.Shift.MattermostUserIDs[user.MattermostUserID] = store.NotEmpty
		joined[user.MattermostUserID] = user
	}

	err := sl.addEventToUsers(joined, NewShiftEvent(rotation, shiftNumber, shift), persist)
	if err != nil {
		return nil, err
	}

	return joined, nil
}
