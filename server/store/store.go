// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-msoffice/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-msoffice/server/utils/kvstore"
)

const (
	UserKeyPrefix = "user_"
)

const OAuth2KeyExpiration = 15 * time.Minute

var ErrNotFound = kvstore.ErrNotFound

type Store interface {
	UserStore
}

type pluginStore struct {
	basicKV kvstore.KVStore
	userKV  kvstore.KVStore
	Logger  bot.Logger
}

func NewPluginStore(api plugin.API, logger bot.Logger) Store {
	basicKV := kvstore.NewPluginStore(api)
	return &pluginStore{
		basicKV: basicKV,
		userKV:  kvstore.NewHashedKeyStore(basicKV, UserKeyPrefix),
		Logger:  logger,
	}
}
