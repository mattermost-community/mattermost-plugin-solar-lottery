// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/config"

import "fmt"

func (api *api) dmUser(user *User, message string) {
	api.Poster.DM(user.MattermostUserID, message)
	api.Debugf("DM bot to %s:\n%s", api.MarkdownUser(user), message)
}

func (api *api) messageWelcomeNewUser(user *User) {
	api.ExpandUser(user)

	// There is the special case when a user uses the plugin for the first time,
	// in which case the actingUser is not yet set. Default to the "subject" user.
	actingUser := api.actingUser
	if actingUser == nil {
		actingUser = user
	}

	api.dmUser(user,
		fmt.Sprintf("### Welcome to Solar Lottery!\n"+
			"%s added you to the Solar Lottery team rotation scheduler. Please use `/%s info` for more information.",
			api.MarkdownUser(actingUser),
			config.CommandTrigger))
}

func (api *api) messageWelcomeToRotation(user *User, rotation *Rotation) {
	api.dmUser(user,
		fmt.Sprintf("### Welcome to %s!\n"+
			"%s added you to %s. Please use `/%s info` for more information.\n"+
			"%s",
			MarkdownRotation(rotation),
			api.MarkdownUser(api.actingUser),
			MarkdownRotation(rotation),
			config.CommandTrigger,
			api.MarkdownRotationBullets(rotation)))
}

func (api *api) messageLeftRotation(user *User, rotation *Rotation) {
	api.dmUser(user,
		fmt.Sprintf("%s removed you from %s.",
			api.MarkdownUser(api.actingUser),
			MarkdownRotation(rotation)))
}

func (api *api) messageAddedSkill(user *User, skillName string, level int) {
	api.ExpandUser(user)
	if level == 0 {
		api.dmUser(user,
			fmt.Sprintf("%s added skill %s, level %s to your profile.\n"+
				"Your current skills are: %s.\n",
				api.MarkdownUser(api.actingUser),
				skillName,
				Level(level),
				MarkdownUserSkills(user)))
	} else {
		api.dmUser(user,
			fmt.Sprintf("%s deleted skill %v from your profile.\n"+
				"Your current skills are: %s.\n",
				api.MarkdownUser(api.actingUser),
				skillName,
				MarkdownUserSkills(user)))
	}
}

func (api *api) messageShiftOpened(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.Users {
		api.dmUser(user,
			fmt.Sprintf("%s opened %s.\n"+
				"Use `/%s shift join -r %s -s %v` if you would like to participate.\n",
				api.MarkdownUser(api.actingUser),
				MarkdownShift(rotation, shiftNumber),
				config.CommandTrigger,
				rotation.Name,
				shiftNumber))
	}
}

func (api *api) messageShiftStarted(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		api.dmUser(user,
			fmt.Sprintf("###### Your %s started!\n"+
				"%s started %s.\n\nTODO runbook URL/channel",
				MarkdownShift(rotation, shiftNumber),
				api.MarkdownUser(api.actingUser),
				MarkdownShift(rotation, shiftNumber)))
	}
}

func (api *api) messageShiftWillStart(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {

		api.dmUser(user,
			fmt.Sprintf("Your %s will start on %s\n\nTODO runbook URL/channel",
				MarkdownShift(rotation, shiftNumber),
				shift.Start))
	}
}

func (api *api) messageShiftFinished(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		api.dmUser(user,
			fmt.Sprintf("###### Done with %s!\n"+
				"%s finished %s. Details:\n%s",
				MarkdownShift(rotation, shiftNumber),
				api.MarkdownUser(api.actingUser),
				MarkdownShift(rotation, shiftNumber),
				api.MarkdownShiftBullets(rotation, shiftNumber, shift)))
	}
}

func (api *api) messageShiftWillFinish(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		api.dmUser(user,
			fmt.Sprintf("Your %s will finish on %s\n\nTODO runbook URL/channel",
				MarkdownShift(rotation, shiftNumber),
				shift.End))
	}
}

func (api *api) messageShiftJoined(joined UserMap, rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	// Notify the previous shift users that new volunteers have been added
	for _, user := range rotation.ShiftUsers(shift) {
		if joined[user.MattermostUserID] != nil {
			continue
		}
		api.dmUser(user,
			fmt.Sprintf("%s added users %s to your %s",
				api.MarkdownUser(api.actingUser),
				api.MarkdownUsers(joined),
				MarkdownShift(rotation, shiftNumber)))
	}

	for _, user := range joined {
		api.dmUser(user,
			fmt.Sprintf("%s joined you into %s",
				api.MarkdownUser(api.actingUser),
				MarkdownShift(rotation, shiftNumber)))
	}
}
