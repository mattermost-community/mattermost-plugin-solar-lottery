// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
)

func (sl *solarLottery) dmUser(user *User, message string) {
	sl.Poster.DM(user.MattermostUserID, message)
	sl.Debugf("DM bot to %s:\n%s", user.Markdown(), message)
}

func (sl *solarLottery) messageWelcomeNewUser(user *User) {
	sl.ExpandUser(user)

	// There is the special case when a user uses the plugin for the first time,
	// in which case the actingUser is not yet set. Default to the "subject" user.
	actingUser := sl.actingUser
	if actingUser == nil {
		actingUser = user
	}

	sl.dmUser(user,
		fmt.Sprintf("### Welcome to Solar Lottery!\n"+
			"%s added you to the Solar Lottery team rotation scheduler. Please use `/%s info` for more information.",
			actingUser.Markdown(),
			config.CommandTrigger))
}

func (sl *solarLottery) messageWelcomeToRotation(user *User, rotation *Rotation) {
	sl.dmUser(user,
		fmt.Sprintf("### Welcome to %s!\n"+
			"%s added you to %s. Please use `/%s info` for more information.\n"+
			"%s",
			rotation.Markdown(),
			sl.actingUser.Markdown(),
			rotation.Markdown(),
			config.CommandTrigger,
			rotation.MarkdownBullets()))
}

func (sl *solarLottery) messageLeftRotation(user *User, rotation *Rotation) {
	sl.dmUser(user,
		fmt.Sprintf("%s removed you from %s.",
			sl.actingUser.Markdown(),
			rotation.Markdown()))
}

func (sl *solarLottery) messageAddedSkill(user *User, skillName string, level int) {
	sl.ExpandUser(user)
	if level == 0 {
		sl.dmUser(user,
			fmt.Sprintf("%s added skill %s, level %s to your profile.\n"+
				"Your current skills are: %s.\n",
				sl.actingUser.Markdown(),
				skillName,
				Level(level),
				user.MarkdownSkills()))
	} else {
		sl.dmUser(user,
			fmt.Sprintf("%s deleted skill %v from your profile.\n"+
				"Your current skills are: %s.\n",
				sl.actingUser.Markdown(),
				skillName,
				user.MarkdownSkills()))
	}
}

func (sl *solarLottery) messageShiftOpened(rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	for _, user := range rotation.Users {
		sl.dmUser(user,
			fmt.Sprintf("%s opened %s.\n"+
				"Use `/%s shift join -r %s -s %v` if you would like to participate.\n",
				sl.actingUser.Markdown(),
				shift.Markdown(),
				config.CommandTrigger,
				rotation.Name,
				shift.ShiftNumber))
	}
}

func (sl *solarLottery) messageShiftStarted(rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		sl.dmUser(user,
			fmt.Sprintf("###### Your %s started!\n"+
				"%s started %s.\n\nTODO runbook URL/channel",
				shift.Markdown(),
				sl.actingUser.Markdown(),
				shift.Markdown()))
	}
}

func (sl *solarLottery) messageShiftWillStart(rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {

		sl.dmUser(user,
			fmt.Sprintf("Your %s will start on %s\n\nTODO runbook URL/channel",
				shift.Markdown(),
				shift.Start))
	}
}

func (sl *solarLottery) messageShiftFinished(rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		sl.dmUser(user,
			fmt.Sprintf("###### Done with %s!\n"+
				"%s finished %s. Details:\n%s",
				shift.Markdown(),
				sl.actingUser.Markdown(),
				shift.Markdown(),
				shift.MarkdownBullets(rotation)))
	}
}

func (sl *solarLottery) messageShiftWillFinish(rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		sl.dmUser(user,
			fmt.Sprintf("Your %s will finish on %s\n\nTODO runbook URL/channel",
				shift.Markdown(),
				shift.End))
	}
}

func (sl *solarLottery) messageShiftJoined(joined UserMap, rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	// Notify the previous shift users that new volunteers have been added
	for _, user := range rotation.ShiftUsers(shift) {
		if joined[user.MattermostUserID] != nil {
			continue
		}
		sl.dmUser(user,
			fmt.Sprintf("%s added users %s to your %s",
				sl.actingUser.Markdown(),
				joined.Markdown(),
				shift.Markdown()))
	}

	for _, user := range joined {
		sl.dmUser(user,
			fmt.Sprintf("%s joined you into %s",
				sl.actingUser.Markdown(),
				shift.Markdown()))
	}
}

func (sl *solarLottery) messageShiftLeft(deleted UserMap, rotation *Rotation, shift *Shift) {
	sl.ExpandRotation(rotation)

	// Notify the previous shift users that users have been deleted from the shift
	for _, user := range rotation.ShiftUsers(shift) {
		if deleted[user.MattermostUserID] != nil {
			continue
		}
		sl.dmUser(user,
			fmt.Sprintf("%s removed users %s from your %s",
				sl.actingUser.Markdown(),
				deleted.Markdown(),
				shift.Markdown()))
	}

	for _, user := range deleted {
		sl.dmUser(user,
			fmt.Sprintf("%s removed you from %s.",
				sl.actingUser.Markdown(),
				shift.Markdown()))
	}
}
