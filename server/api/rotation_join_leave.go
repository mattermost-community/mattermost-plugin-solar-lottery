// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) JoinRotation(mattermostUsernames string, rotation *Rotation, starting time.Time) (UserMap, error) {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.AddRotationUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"RotationID":          rotation.RotationID,
		"MattermostUsernames": mattermostUsernames,
		"Starting":            starting.Format(DateFormat),
	})

	// -1 is acceptable
	shiftNumber, _ := rotation.ShiftNumberForTime(starting)
	added := UserMap{}
	for _, user := range api.users {
		if len(rotation.MattermostUserIDs[user.MattermostUserID]) != 0 {
			logger.Debugf("%s is already in rotation %s.",
				added.MarkdownWithSkills(), rotation.Markdown())
			continue
		}

		// A new person may be given some slack - setting LastShiftNumber in the
		// future guarantees they won't be selected until then.
		user.LastServed[rotation.RotationID] = shiftNumber

		user, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return added, err
		}

		if rotation.MattermostUserIDs == nil {
			rotation.MattermostUserIDs = store.IDMap{}
		}

		rotation.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID
		api.messageWelcomeToRotation(user, rotation)
		added[user.MattermostUserID] = user
	}

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return added, errors.WithMessagef(err, "failed to store rotation %s", rotation.RotationID)
	}
	logger.Infof("%s added %s to %s.",
		api.actingUser.Markdown(), added.MarkdownWithSkills(), rotation.Markdown())
	return added, nil
}

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
			logger.Debugf("%s is not found in rotation %s", user.Markdown(), rotation.Markdown())
			continue
		}

		delete(user.LastServed, rotation.RotationID)
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

	logger.Infof("%s removed from %s.", deleted.Markdown(), rotation.Markdown())
	return deleted, nil
}
