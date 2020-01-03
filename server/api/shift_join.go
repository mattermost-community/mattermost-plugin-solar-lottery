// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrUserAlreadyInShift = errors.New("user is already in shift")

func (api *api) JoinShift(mattermostUsernames string, rotation *Rotation, shiftNumber int) (*Shift, UserMap, error) {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return nil, nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.VolunteerUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"RotationID":          rotation.RotationID,
		"ShiftNumber":         shiftNumber,
	})

	shift, joined, err := api.joinShift(api.users, rotation, shiftNumber, true)
	if err != nil {
		return nil, nil, err
	}
	api.messageShiftJoined(joined, rotation, shiftNumber, shift)
	logger.Infof("%s volunteered %s to %s.",
		api.MarkdownUser(api.actingUser), api.MarkdownUsersWithSkills(joined), MarkdownShift(rotation, shiftNumber))
	return shift, joined, nil
}

func (api *api) joinShift(users UserMap, rotation *Rotation, shiftNumber int, persist bool) (*Shift, UserMap, error) {
	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return nil, nil, errors.Errorf("failed to load shift %v for rotation %s", shiftNumber, rotation.RotationID)
	}
	if shift.Status != store.ShiftStatusOpen {
		return nil, nil, errors.Errorf("can't volunteer for a status which is %s, must be Open", shift.Status)
	}

	joined := UserMap{}
	for _, user := range users {
		if shift.Shift.MattermostUserIDs[user.MattermostUserID] != "" {
			continue
		}
		if len(shift.MattermostUserIDs) >= rotation.Size {
			return nil, nil, errors.Errorf("rotation size %v exceeded", rotation.Size)
		}
		shift.Shift.MattermostUserIDs[user.MattermostUserID] = store.NotEmpty
		joined[user.MattermostUserID] = user
	}

	if persist {
		err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
		if err != nil {
			return nil, nil, err
		}
	}

	err = api.addEventToUsers(joined, NewShiftEvent(rotation, shiftNumber, shift), persist)
	if err != nil {
		return nil, nil, err
	}

	return shift, joined, nil
}
