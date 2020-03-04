// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

func (sl *sl) AddToCalendar(mattermostUserIDs *types.IDSet, u *Unavailable) (*Users, error) {
	var users *Users
	err := sl.Setup(
		pushLogger("AddToCalendar", bot.LogContext{ctxUnavailable: u}),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	for _, user := range users.AsArray() {
		sl.addUserUnavailable(user, u)
	}

	sl.Infof("%s added event %s to %s.",
		sl.actingUser.Markdown(), sl.actingUser.MarkdownUnavailable(u), users.MarkdownWithSkills())
	return users, nil
}

func (sl *sl) addUserUnavailable(user *User, u *Unavailable) error {
	user.AddUnavailable(u)
	err := sl.storeUser(user)
	if err != nil {
		return errors.Wrapf(err, "user: %s", user.Markdown())
	}
	return nil
}

func (sl *sl) ClearCalendar(mattermostUserIDs *types.IDSet, interval types.Interval) (*Users, error) {
	var users *Users
	err := sl.Setup(
		pushLogger("CkearCalendar", bot.LogContext{ctxInterval: interval}),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	for _, user := range users.AsArray() {
		_ = user.findUnavailable(interval, true)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to remove unavailable for %v", interval)
		}

		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to update user %s", user.Markdown())
		}
	}

	sl.Infof("%s deleted events %v from users %s.",
		sl.actingUser.Markdown(), interval, users.MarkdownWithSkills())
	return users, nil
}

func (sl *sl) clearUserUnavailable(user *User, interval types.Interval) error {
	_ = user.findUnavailable(interval, true)
	_, err := sl.storeUserWelcomeNew(user)
	if err != nil {
		return errors.Wrapf(err, "user %s", user.Markdown())
	}
	return nil
}
