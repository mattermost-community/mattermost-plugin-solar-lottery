// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

const (
	UserKeyPrefix      = "user_"
	SkillsKeyPrefix    = "skills_"
	RotationsKeyPrefix = "rotations_"
)

const OAuth2KeyExpiration = 15 * time.Minute

var ErrNotFound = kvstore.ErrNotFound

type Store interface {
	UserStore
	SkillsStore
	RotationsStore
}

type pluginStore struct {
	basicKV     kvstore.KVStore
	userKV      kvstore.KVStore
	skillsKV    kvstore.KVStore
	rotationsKV kvstore.KVStore
	Logger      bot.Logger
}

func NewPluginStore(api plugin.API, logger bot.Logger) Store {
	basicKV := kvstore.NewPluginStore(api)
	return &pluginStore{
		basicKV:     basicKV,
		userKV:      kvstore.NewHashedKeyStore(basicKV, UserKeyPrefix),
		skillsKV:    kvstore.NewHashedKeyStore(basicKV, SkillsKeyPrefix),
		rotationsKV: kvstore.NewHashedKeyStore(basicKV, RotationsKeyPrefix),
		Logger:      logger,
	}
}
