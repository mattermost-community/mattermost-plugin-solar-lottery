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
	bot.Logger
	User
	Skills
	Rotations
}

type PluginAPI interface {
	GetUser(string) (*model.User, error)
	GetUserByUsername(string) (*model.User, error)
	IsPluginAdmin(mattermostUserID string) (bool, error)
}

// Dependencies contains all API dependencies
type Dependencies struct {
	Logger         bot.Logger
	PluginAPI      PluginAPI
	Poster         bot.Poster
	RotationsStore store.RotationsStore
	ShiftStore     store.ShiftStore
	SkillsStore    store.SkillsStore
	UserStore      store.UserStore
}

type Config struct {
	*config.Config
	*Dependencies
}

type api struct {
	bot.Logger
	Config
	mattermostUserID string
	user             *store.User
	skills           []string
	rotations        map[string]*store.Rotation
}

func New(apiConfig Config, mattermostUserID string) API {
	return &api{
		Logger: apiConfig.Logger.With(bot.LogContext{
			"MattermostUserID": mattermostUserID,
		}),
		Config:           apiConfig,
		mattermostUserID: mattermostUserID,
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
