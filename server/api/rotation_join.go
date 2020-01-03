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
				api.MarkdownUsersWithSkills(added), MarkdownRotation(rotation))
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
		api.MarkdownUser(api.actingUser), api.MarkdownUsersWithSkills(added), MarkdownRotation(rotation))
	return added, nil
}
