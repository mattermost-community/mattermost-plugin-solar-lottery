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
			api.by(user), config.CommandTrigger, MarkdownRotationWithDetails(rotation))
}

func (api *api) messageLeftRotation(user *User, rotation *Rotation) {
	api.Poster.DM(user.MattermostUserID,
		"You have been removed from the rotation %s%s.", MarkdownRotation(rotation), api.by(user))
}

func (api *api) messageAddedSkill(user *User, skillName string, level int) {
	api.ExpandUser(user)
	if level == 0 {
		api.Poster.DM(user.MattermostUserID,
			"Skill %v, level %v was added to your profile%s.\n"+
				"Your current skills are: %s\n",
			skillName, LevelToString(level), api.by(user), MarkdownUserSkills(user))
	} else {
		api.Poster.DM(user.MattermostUserID,
			"Skill %v was deleted from your profile%s.\n"+
				"Your current skills are: %s\n",
			skillName, api.by(user), MarkdownUserSkills(user))
	}
}

func (api *api) by(forUser *User) string {
	if forUser.MattermostUserID == api.actingMattermostUserID {
		return ""
	}
	api.ExpandUser(api.actingUser)
	return " by " + MarkdownUser(api.actingUser)
}
