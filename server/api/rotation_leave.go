// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) LeaveRotation(mattermostUsernames string, rotation *Rotation) (UserMap, error) {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.DeleteUsersFromRotation",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"RotationID":          rotation.RotationID,
	})

	deleted := UserMap{}
	for _, user := range api.users {
		_, ok := rotation.MattermostUserIDs[user.MattermostUserID]
		if !ok {
			logger.Debugf("%s is not found in rotation %s", MarkdownUser(user), MarkdownRotation(rotation))
			continue
		}

		delete(user.NextRotationShift, rotation.RotationID)
		_, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return deleted, err
		}
		delete(rotation.MattermostUserIDs, user.MattermostUserID)
		if len(rotation.Users) > 0 {
			delete(rotation.Users, user.MattermostUserID)
		}
		api.messageLeftRotation(user, rotation)
		deleted[user.MattermostUserID] = user
	}

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return deleted, err
	}

	logger.Infof("%s removed from %s.", MarkdownUserMap(deleted), MarkdownRotation(rotation))
	return deleted, nil
}

func (api *api) deleteUsersFromRotation(rotation *Rotation, users UserMap) error {
	for _, user := range users {
		_, ok := rotation.MattermostUserIDs[user.MattermostUserID]
		if !ok {
			return errors.Errorf("%s is not found in rotation %s", MarkdownUser(user), MarkdownRotation(rotation))
		}

		delete(user.NextRotationShift, rotation.RotationID)
		_, err := api.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
		delete(rotation.MattermostUserIDs, user.MattermostUserID)
		api.messageLeftRotation(user, rotation)
		api.Logger.Debugf("removed %s from %s.", MarkdownUser(user), MarkdownRotation(rotation))
		return nil
	}

	return api.RotationStore.StoreRotation(rotation.Rotation)
}
