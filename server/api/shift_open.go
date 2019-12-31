// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrShiftAlreadyExists = errors.New("shift already exists")

func (api *api) OpenShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.OpenShift",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"ShiftNumber":    shiftNumber,
	})

	_, err = api.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	if err != store.ErrNotFound {
		if err != nil {
			return nil, err
		}
		return nil, ErrShiftAlreadyExists
	}

	shift, err := rotation.makeShift(shiftNumber, nil)
	if err != nil {
		return nil, err
	}
	shift.ShiftStatus = store.ShiftStatusOpen

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return nil, err
	}

	api.messageShiftOpened(rotation, shiftNumber, shift)
	logger.Infof("%s opened shift %s in %s.", MarkdownUser(api.actingUser), MarkdownShift(shiftNumber, shift), MarkdownRotation(rotation))
	return shift, nil
}
