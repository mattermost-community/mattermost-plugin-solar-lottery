// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
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

// Dependencies contains all API dependencies
type Dependencies struct {
	UserStore         store.UserStore
	RotationsStore    store.RotationsStore
	SkillsStore       store.SkillsStore
	ShiftStore        store.ShiftStore
	Logger            bot.Logger
	Poster            bot.Poster
	IsAuthorizedAdmin func(userID string) (bool, error)
}

type Config struct {
	*Dependencies
	*config.Config
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
