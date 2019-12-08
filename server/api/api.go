// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type API interface {
	User
	Skills
}

// Dependencies contains all API dependencies
type Dependencies struct {
	UserStore         store.UserStore
	RotationStore     store.RotationStore
	SkillsStore       store.SkillsStore
	Logger            bot.Logger
	Poster            bot.Poster
	IsAuthorizedAdmin func(userID string) (bool, error)
}

type Config struct {
	*Dependencies
	*config.Config
}

type api struct {
	Config
	mattermostUserID string
	user             *store.User
}

func New(apiConfig Config, mattermostUserID string) API {
	return &api{
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

func withUser(api *api) error {
	if api.user != nil {
		return nil
	}

	user, err := api.UserStore.LoadUser(api.mattermostUserID)
	if err != nil {
		return err
	}

	api.user = user
	return nil
}
