// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) UpdateRotation(rotation *Rotation, updatef func(*Rotation) error) error {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.UpdateRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})

	err = updatef(rotation)
	if err != nil {
		return err
	}

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}

	logger.Infof("%s updated rotation %s.", api.MarkdownUser(api.actingUser), MarkdownRotation(rotation))
	return nil
}
