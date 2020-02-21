// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type PluginAPI interface {
	GetMattermostUser(mattermostUserID string) (*model.User, error)
	GetMattermostUserByUsername(mattermostUsername string) (*model.User, error)
	IsPluginAdmin(mattermostUserID string) (bool, error)
	Clean() error
	GetBotUserID() string
}

type Service struct {
	PluginAPI
	Config config.Service

	// Autofillers map[string]Autofiller
	Logger bot.Logger
	Poster bot.Poster
	Store  kvstore.Store
}

func (s *Service) ActingAs(mattermostUserID string) SL {
	return &sl{
		Service:                s,
		conf:                   s.Config.Get(),
		actingMattermostUserID: mattermostUserID,
		Logger: s.Logger.With(bot.LogContext{
			"ActingUserID": mattermostUserID,
		}),
	}
}

func (s *Service) Clean() error {
	return s.PluginAPI.Clean()
}
