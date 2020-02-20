// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Tasker interface {
	// FillTask(*Rotation, *Task) (added UserMap, err error)
	// PostTask(*Rotation, *Task) error
}

type PluginAPI interface {
	GetMattermostUser(mattermostUserID string) (*model.User, error)
	GetMattermostUserByUsername(mattermostUsername string) (*model.User, error)
	IsPluginAdmin(mattermostUserID string) (bool, error)
	Clean() error
}

type SL interface {
	Tasker
	// Autopilot

	Rotations
	Calendar
	// Tasks
	Skills
	Users

	PluginAPI
}

// Dependencies contains all API dependencies
type Service struct {
	*config.Config
	// Autofillers map[string]Autofiller
	PluginAPI
	Logger bot.Logger
	Poster bot.Poster
	Store  kvstore.Store
}

type sl struct {
	*Service
	bot.Logger

	// set by `sl.New`
	actingMattermostUserID string

	// use withActingUser or withActingUserExpanded to initialize.
	actingUser *User

	// use withKnownSkills to initialize.
	knownSkills *types.Set

	// use withActiveRotations or withRotation(rotationID) to initialize, not expanded by default.
	activeRotations *types.Set

	// // use withMattermostUsers(usernames) or withUsers(mattermostUserIDs) to initialize, not expanded by default.
	// users UserMap
}

func (s *Service) ActingAs(mattermostUserID string) SL {
	return &sl{
		Service:                s,
		actingMattermostUserID: mattermostUserID,
		Logger: s.Logger.With(bot.LogContext{
			"ActingUserID": mattermostUserID,
		}),
	}
}

func (s *Service) Clean() error {
	return s.PluginAPI.Clean()
}
