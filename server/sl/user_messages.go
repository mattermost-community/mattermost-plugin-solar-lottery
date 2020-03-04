// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
)

func (sl *sl) dmUser(user *User, message string) {
	sl.Poster.DM(string(user.MattermostUserID), message)
	sl.Debugf("DM bot to %s:\n%s", user.Markdown(), message)
}

func (sl *sl) messageWelcomeNewUser(user *User) {
	sl.expandUser(user)

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
			constants.CommandTrigger))
}

func (sl *sl) messageWelcomeToRotation(user *User, rotation *Rotation) {
	sl.dmUser(user,
		fmt.Sprintf("### Welcome to %s!\n"+
			"%s added you to %s. Please use `/%s info` for more information.\n"+
			"%s",
			rotation.Markdown(),
			sl.actingUser.Markdown(),
			rotation.Markdown(),
			constants.CommandTrigger,
			rotation.MarkdownBullets()))
}

func (sl *sl) messageLeftRotation(user *User, rotation *Rotation) {
	sl.dmUser(user,
		fmt.Sprintf("%s removed you from %s.",
			sl.actingUser.Markdown(),
			rotation.Markdown()))
}

func (sl *sl) messageAddedSkill(user *User, skillName string, level int) {
	sl.expandUser(user)
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

func (sl *sl) messageNewTask(rotation *Rotation, t *Task) {
	sl.expandRotationUsers(rotation)

	for _, user := range rotation.users.AsArray() {
		sl.dmUser(user,
			fmt.Sprintf("%s opened %s.\n"+
				"Use `TODO` if you would like to participate.\n",
				sl.actingUser.Markdown(),
				t.Markdown()))
	}
}

func (sl *sl) messageTaskStarted(rotation *Rotation, t *Task) {
	sl.expandRotationUsers(rotation)

	for _, user := range rotation.FindUsers(t.MattermostUserIDs) {
		sl.dmUser(user,
			fmt.Sprintf("###### Your %s started!\n"+
				"%s started %s.\n\nTODO runbook URL/channel",
				t.Markdown(),
				sl.actingUser.Markdown(),
				t.Markdown()))
	}
}

func (sl *sl) messageTaskWillStart(rotation *Rotation, t *Task) {
	sl.expandRotationUsers(rotation)

	for _, user := range rotation.FindUsers(t.MattermostUserIDs) {
		sl.dmUser(user,
			fmt.Sprintf("Your %s will start on TODO\n\nTODO runbook URL/channel",
				t.Markdown()))
	}
}

func (sl *sl) messageTaskFinished(rotation *Rotation, t *Task) {
	sl.expandRotationUsers(rotation)

	for _, user := range rotation.FindUsers(t.MattermostUserIDs) {
		sl.dmUser(user,
			fmt.Sprintf("###### Done with %s!\n"+
				"%s finished %s. Details:\n%s",
				t.Markdown(),
				sl.actingUser.Markdown(),
				t.Markdown(),
				t.MarkdownBullets(rotation)))
	}
}

func (sl *sl) messageTaskWillFinish(rotation *Rotation, t *Task) {
	sl.expandRotationUsers(rotation)

	for _, user := range rotation.FindUsers(t.MattermostUserIDs) {
		sl.dmUser(user,
			fmt.Sprintf("Your %s will finish on TODO\n\nTODO runbook URL/channel",
				t.Markdown()))
	}
}

func (sl *sl) messageAddedUsersToTask(added Users, rotation *Rotation, t *Task) {
	sl.expandRotationUsers(rotation)

	// Notify the previous shift users that new volunteers have been added
	for _, user := range rotation.FindUsers(t.MattermostUserIDs) {
		if !added.Contains(user.MattermostUserID) {
			continue
		}
		sl.dmUser(user,
			fmt.Sprintf("%s added users %s to your %s",
				sl.actingUser.Markdown(),
				added.Markdown(),
				t.Markdown()))
	}

	for _, user := range added.AsArray() {
		sl.dmUser(user,
			fmt.Sprintf("%s assigned you to %s",
				sl.actingUser.Markdown(),
				t.Markdown()))
	}
}
