// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

type Calendar interface {
	AddToCalendar(users UserMap, u *Unavailable) error
	ClearCalendar(users UserMap, interval types.Interval) error
}

func (sl *sl) AddToCalendar(users UserMap, u *Unavailable) error {
	err := sl.Setup(pushLogger("AddToCalendar",
		bot.LogContext{
			ctxUsernames:   users.String(),
			ctxUnavailable: u,
		}))
	if err != nil {
		return err
	}
	defer sl.popLogger()

	for _, user := range users {
		user.AddUnavailable(u)
		err = sl.storeUser(user)
		if err != nil {
			return errors.WithMessagef(err, "failed to update user %s", user.Markdown())
		}
	}

	sl.Infof("%s added event %s to %s.",
		sl.actingUser.Markdown(), sl.actingUser.MarkdownUnavailable(u), users.MarkdownWithSkills())
	return nil
}

func (sl *sl) ClearCalendar(users UserMap, interval types.Interval) error {
	err := sl.Setup(pushLogger("CkearCalendar",
		bot.LogContext{
			ctxUsernames: users.String(),
			ctxInterval:  interval,
		}))
	if err != nil {
		return err
	}
	defer sl.popLogger()

	for _, user := range users {
		_ = user.findUnavailable(interval, true)
		if err != nil {
			return errors.WithMessagef(err, "failed to remove unavailable for %v", interval)
		}

		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return errors.WithMessagef(err, "failed to update user %s", user.Markdown())
		}
	}

	sl.Infof("%s deleted events %v from users %s.",
		sl.actingUser.Markdown(), interval, users.MarkdownWithSkills())
	return nil
}
