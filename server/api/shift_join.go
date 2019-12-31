// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrUserAlreadyInShift = errors.New("user is already in shift")

func (api *api) JoinShift(mattermostUsernames string, rotation *Rotation, shiftNumber int) error {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.VolunteerUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"RotationID":          rotation.RotationID,
		"ShiftNumber":         shiftNumber,
	})

	shift, err := api.loadShift(rotation, shiftNumber)
	if err != nil {
		return errors.Errorf("failed to load shift %v for rotation %s", shiftNumber, rotation.RotationID)
	}
	if shift.Status != store.ShiftStatusOpen {
		return errors.Errorf("can't volunteer for a status which is %s, must be Open", shift.Status)
	}

	volunteered := UserMap{}
	for _, user := range api.users {
		// TODO error if the shift has no vacancies? Or allow volunteering above the limit, and let the committer choose?
		if shift.Shift.MattermostUserIDs[user.MattermostUserID] != "" {
			return ErrUserAlreadyInShift
		}
		shift.Shift.MattermostUserIDs[user.MattermostUserID] = store.NotEmpty
		shift.Users[user.MattermostUserID] = user
		volunteered[user.MattermostUserID] = user
	}

	err = api.ShiftStore.StoreShift(rotation.RotationID, shiftNumber, shift.Shift)
	if err != nil {
		return err
	}

	api.messageShiftVolunteers(volunteered, rotation, shiftNumber, shift)
	logger.Infof("%s volunteered %s to %s.",
		MarkdownUser(api.actingUser), MarkdownUserMapWithSkills(volunteered), MarkdownShift(rotation, shiftNumber, shift))
	return nil
}
