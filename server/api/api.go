// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type API interface {
	Users
	Skills
	Rotations
	Forecaster

	bot.Logger
	PluginAPI
}

type PluginAPI interface {
	GetMattermostUser(mattermostUserID string) (*model.User, error)
	GetMattermostUserByUsername(mattermostUsername string) (*model.User, error)
	IsPluginAdmin(mattermostUserID string) (bool, error)
}

// Dependencies contains all API dependencies
type Dependencies struct {
	PluginAPI
	Logger        bot.Logger
	Poster        bot.Poster
	RotationStore store.RotationStore
	ShiftStore    store.ShiftStore
	SkillsStore   store.SkillsStore
	UserStore     store.UserStore
}

type Config struct {
	*config.Config
	*Dependencies
}

type api struct {
	bot.Logger
	Config

	// set by `api.New`
	actingMattermostUserID string

	// use withActingUser or withActingUserExpanded to initialize.
	actingUser *User

	// use withKnownSkills to initialize.
	knownSkills store.IDMap

	// use withKnownRotations or withRotation(rotationID) to initialize, not expanded by default.
	knownRotations store.IDMap

	// use withRotation(rotationID) to initialize, not expanded by default.
	// rotation *Rotation

	// use withMattermostUsers(usernames) or withUsers(mattermostUserIDs) to initialize, not expanded by default.
	users UserMap
}

func New(apiConfig Config, mattermostUserID string) API {
	return &api{
		Logger: apiConfig.Logger.With(bot.LogContext{
			"MattermostUserID": mattermostUserID,
		}),
		Config:                 apiConfig,
		actingMattermostUserID: mattermostUserID,
	}
}

type filterf func(*api) error

func (api *api) Filter(filters ...filterf) error {
	for _, filter := range filters {
		err := filter(api)
		if err != nil {
			return err
		}
	}
	return nil
}
