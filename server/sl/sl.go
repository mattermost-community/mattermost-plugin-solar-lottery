// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Tasker interface {
	// FillTask(*Rotation, *Task) (added UserMap, err error)
	// PostTask(*Rotation, *Task) error
}

type SL interface {
	Calendar
	Rotations
	Skills
	Tasker
	Users

	PluginAPI
	bot.Logger
}

type sl struct {
	*Service
	*config.Config
	bot.Logger

	// set by Service.ActingAs.
	actingMattermostUserID string

	// set by withActingUser or withActingUserExpanded.
	actingUser *User

	// Common indices (set by withXXX).
	knownSkills     *types.Set
	activeRotations *types.Set
}
