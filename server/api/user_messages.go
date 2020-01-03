// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/config"

func (api *api) messageWelcomeNewUser(user *User) {
	if user.PluginVersion != "" {
		return
	}

	api.ExpandUser(user)
	api.Poster.DM(user.MattermostUserID,
		"###### Welcome to Solar Lottery!\n"+
			"You have been added to the Solar Lottery team rotation scheduler%s. Please use `%s help` for more information."+
			api.by(user), config.CommandTrigger)
}

func (api *api) messageWelcomeToRotation(user *User, rotation *Rotation) {
	api.Poster.DM(user.MattermostUserID,
		"###### Welcome to rotation %s!\n"+
			"You have been added%s. Please use `%s help` for more information.\n"+
			"%s"+
			api.by(user), config.CommandTrigger, api.MarkdownRotationWithDetails(rotation))
}

func (api *api) messageLeftRotation(user *User, rotation *Rotation) {
	api.Poster.DM(user.MattermostUserID,
		"You have been removed from the rotation %s%s.", MarkdownRotation(rotation), api.by(user))
}

func (api *api) messageAddedSkill(user *User, skillName string, level int) {
	api.ExpandUser(user)
	if level == 0 {
		api.Poster.DM(user.MattermostUserID,
			"Skill %s, level %s was added to your profile%s.\n"+
				"Your current skills are: %s.\n",
			skillName, Level(level), api.by(user), MarkdownUserSkills(user))
	} else {
		api.Poster.DM(user.MattermostUserID,
			"Skill %v was deleted from your profile%s.\n"+
				"Your current skills are: %s.\n",
			skillName, api.by(user), MarkdownUserSkills(user))
	}
}

func (api *api) messageShiftOpened(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.Users {
		api.Poster.DM(user.MattermostUserID,
			"%s opened%s.\n"+
				"Use `/%s shift join -r %s -s %v` if you would like to participate.\n",
			MarkdownShift(rotation, shiftNumber), api.by(user),
			config.CommandTrigger, rotation.Name, shiftNumber)
	}
}

func (api *api) messageShiftStarted(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		api.Poster.DM(user.MattermostUserID,
			"###### Welcome to %s!\n"+
				"%s started%s.\n\nTODO runbook URL/channel",
			MarkdownShift(rotation, shiftNumber),
			MarkdownShift(rotation, shiftNumber), api.by(user))
	}
}

func (api *api) messageShiftWillStart(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {

		api.Poster.DM(user.MattermostUserID,
			"Your %s will start on %s\n\nTODO runbook URL/channel",
			MarkdownShift(rotation, shiftNumber),
			shift.Start)
	}
}

func (api *api) messageShiftFinished(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		api.Poster.DM(user.MattermostUserID,
			"###### Done with your shift in %s!\n"+
				"Your shift in %s is now finished%s. Details:\n%s",
			MarkdownRotation(rotation),
			MarkdownRotation(rotation), api.by(user), MarkdownShift(rotation, shiftNumber))
	}
}

func (api *api) messageShiftWillFinish(rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	for _, user := range rotation.ShiftUsers(shift) {
		api.Poster.DM(user.MattermostUserID,
			"Your %s will finish on %s\n\nTODO runbook URL/channel",
			MarkdownShift(rotation, shiftNumber),
			shift.End)
	}
}

func (api *api) messageShiftJoined(joined UserMap, rotation *Rotation, shiftNumber int, shift *Shift) {
	api.ExpandRotation(rotation)

	// Notify the previous shift users that new volunteers have been added
	for _, user := range rotation.ShiftUsers(shift) {
		if joined[user.MattermostUserID] != nil {
			continue
		}
		api.Poster.DM(user.MattermostUserID,
			"New users %s added your %s%s",
			api.MarkdownUsers(joined), MarkdownShift(rotation, shiftNumber), api.by(user))
	}

	for _, user := range joined {
		api.Poster.DM(user.MattermostUserID,
			"You volunteered for shift %s%s",
			MarkdownShift(rotation, shiftNumber), api.by(user))
	}
}

func (api *api) by(forUser *User) string {
	if forUser.MattermostUserID == api.actingMattermostUserID {
		return ""
	}
	api.ExpandUser(api.actingUser)
	return " by " + api.MarkdownUser(api.actingUser)
}
