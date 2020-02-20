// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"

type Calendar interface {
	AddToCalendar(users UserMap, u *Unavailable) error
	ClearCalendar(users UserMap, interval types.Interval) error
}

func (sl *sl) AddToCalendar(users UserMap, u *Unavailable) error {
	// err := sl.Filter(
	// 	withActingUserExpanded,
	// 	withMattermostUsersExpanded(users),
	// )
	// if err != nil {
	// 	return err
	// }
	// logger := sl.Logger.Timed().With(bot.LogContext{
	// 	"Location":            "sl.AddSkillToUsers",
	// 	"ActingUsername":      sl.actingUser.MattermostUsername(),
	// 	"MattermostUsernames": users.String(),
	// 	"Event":               event,
	// })

	// err = sl.addEventToUsers(sl.users, event, true)
	// if err != nil {
	// 	return err
	// }

	// logger.Infof("%s added event %s to %s.",
	// 	sl.actingUser.Markdown(), event.Markdown(), sl.users.MarkdownWithSkills())
	return nil
}

func (sl *sl) ClearCalendar(users UserMap, interval types.Interval) error {
	// err := sl.Filter(
	// 	withActingUserExpanded,
	// 	withMattermostUsersExpanded(mattermostUsernames),
	// )
	// if err != nil {
	// 	return err
	// }
	// logger := sl.Logger.Timed().With(bot.LogContext{
	// 	"Location":            "sl.AddSkillToUsers",
	// 	"ActingUsername":      sl.actingUser.MattermostUsername(),
	// 	"MattermostUsernames": mattermostUsernames,
	// 	"StartDate":           startDate,
	// 	"EndDate":             endDate,
	// })

	// for _, user := range sl.users {
	// 	intervalStart, intervalEnd, err := ParseDatePair(startDate, endDate)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	_, err = user.OverlapEvents(intervalStart, intervalEnd, true)
	// 	if err != nil {
	// 		return errors.WithMessagef(err, "failed to remove events from %s to %s", startDate, endDate)
	// 	}

	// 	_, err = sl.storeUserWelcomeNew(user)
	// 	if err != nil {
	// 		return errors.WithMessagef(err, "failed to update user %s", user.Markdown())
	// 	}
	// }

	// logger.Infof("%s deleted events from %s to %s from users %s.",
	// 	sl.actingUser.Markdown(), startDate, endDate, sl.users.MarkdownWithSkills())
	return nil
}
