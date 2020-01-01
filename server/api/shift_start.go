// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) StartShift(rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.StartShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return err
	}

	shift.Status = store.ShiftStatusStarted

	for _, user := range shift.Users {
		user.NextRotationShift[rotation.RotationID] = shiftNumber + rotation.Grace
		_, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
	}

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return err
	}

	api.messageShiftStarted(rotation, shiftNumber, shift)
	logger.Infof("%s started %s.", MarkdownUser(api.actingUser), MarkdownShift(rotation, shiftNumber, shift))
	return nil
}
