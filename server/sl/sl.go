// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const (
	ctxAPI            = "API"
	ctxActingUserID   = "ActingUserID"
	ctxActingUsername = "ActingUsername"
	ctxRotationID     = "RotationID"
	ctxUsers          = "Users"
	ctxSkill          = "Skill"
	ctxSkillLevel     = "SkillLevel"
	ctxStarting       = "Starting"
	ctxUsernames      = "Usernames"
	ctxUnavailable    = "Unavailable"
	ctxInterval       = "Interval"
)

type SL interface {
	Calendar
	Rotations
	Skills
	Issues
	Users

	PluginAPI
	bot.Logger

	ActingUser() (*User, error)
	Config() *config.Config
}

type sl struct {
	*Service
	bot.Logger

	conf *config.Config

	// set by Service.ActingAs.
	actingMattermostUserID types.ID

	// set by withActingUser or withActingUserExpanded.
	actingUser *User

	// Common indices (set by withXXX).
	knownSkills     *types.IDIndex
	activeRotations *types.IDIndex
	loggers         []bot.Logger
}

func (sl *sl) Config() *config.Config {
	return sl.conf
}
