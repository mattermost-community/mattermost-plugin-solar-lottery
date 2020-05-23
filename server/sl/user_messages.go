// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
)

func (sl *sl) dmUserWelcomeToSolarLottery(user *User) {
	sl.expandUser(user)
	// There is the special case when a user uses the plugin for the first time,
	// in which case the actingUser is not yet set. Default to the "subject" user.
	actingUser := sl.actingUser
	if actingUser == nil || actingUser.mattermostUser == nil {
		actingUser = user
	}

	sl.dmUser(user,
		fmt.Sprintf("### Welcome to Solar Lottery!\n"+
			"%s added you to the Solar Lottery team rotation scheduler. Please use `/%s info` for more information.",
			actingUser.Markdown(),
			constants.CommandTrigger))
}

func (sl *sl) dmUserWelcomeToRotation(user *User, rotation *Rotation) {
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

func (sl *sl) dmUserLeftRotation(user *User, rotation *Rotation) {
	sl.dmUser(user,
		fmt.Sprintf("%s removed you from %s.",
			sl.actingUser.Markdown(),
			rotation.Markdown()))
}

func (sl *sl) dmUserChangedSkill(user *User, skillName string, level int) {
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

func (sl *sl) dmUserTaskPending(user *User, task *Task) {
	sl.dmUser(user,
		fmt.Sprintf("%s opened a new pending task %s.\n"+
			"Use `TODO` if you would like to participate.\n",
			sl.actingUser.Markdown(),
			task.Markdown()))
}

func (sl *sl) dmUserTaskStarted(user *User, task *Task) {
	sl.dmUser(user,
		fmt.Sprintf("###### Your %s started!\n"+
			"%s started %s.\n\nTODO runbook URL/channel",
			task.Markdown(),
			sl.actingUser.Markdown(),
			task.Markdown()))
}

func (sl *sl) dmUserTaskScheduled(user *User, t *Task) {
	sl.dmUser(user,
		fmt.Sprintf("###### You have been scheduled for %s.\n"+
			"%s scheduled %s.\n\nTODO runbook/info URL/channel",
			t.Markdown(),
			sl.actingUser.Markdown(),
			t.Markdown()))
}

func (sl *sl) dmUserTaskWillStart(user *User, t *Task) {
	sl.dmUser(user,
		fmt.Sprintf("###### Your task %s will start TODO-when.\n"+
			"TODO runbook/info URL/channel",
			t.Markdown()))
}

func (sl *sl) dmUserTaskFinished(user *User, task *Task) {
	sl.dmUser(user,
		fmt.Sprintf("###### Your %s finished!\n"+
			"%s finished %s.\n\nTODO runbook URL/channel",
			task.Markdown(),
			sl.actingUser.Markdown(),
			task.Markdown()))
}

func (sl *sl) dmUserTaskWillFinish(user *User, task *Task) {
	sl.dmUser(user,
		fmt.Sprintf("###### Your task %s will finish TODO-when.\n"+
			"TODO runbook/info URL/channel",
			task.Markdown()))
}

func (sl *sl) dmUserAssignedTask(user *User, task *Task) {
	if task.State == TaskStatePending {
		return
	}

	sl.dmUser(user,
		fmt.Sprintf("%s assigned you to %s, which is %s",
			sl.actingUser.Markdown(),
			task.Markdown(),
			task.State))
}

func (sl *sl) dmUser(user *User, message string) {
	sl.Poster.DM(string(user.MattermostUserID), message)
	sl.Debugf("DM bot to %s:\n%s", user.Markdown(), message)
}

func (sl *sl) announceRotationUsers(r *Rotation, dm func(*User, *Rotation)) {
	sl.expandRotationUsers(r)
	for _, user := range r.Users.AsArray() {
		dm(user, r)
	}
}

func (sl *sl) announceTaskUsers(task *Task, dm func(*User, *Task)) {
	for _, user := range task.Users.AsArray() {
		dm(user, task)
	}
}
