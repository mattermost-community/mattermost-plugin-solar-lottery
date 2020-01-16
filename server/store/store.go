// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

const (
	UserKeyPrefix     = "user_"
	RotationKeyPrefix = "rotation_"
	ShiftKeyPrefix    = "shift_"

	KnownSkillsKey    = "index_skills"
	KnownRotationsKey = "index_rotations"
)

const OAuth2KeyExpiration = 15 * time.Minute

var ErrNotFound = kvstore.ErrNotFound

type Store interface {
	UserStore
	SkillsStore
	RotationStore
	ShiftStore
}

type pluginStore struct {
	basicKV    kvstore.KVStore
	userKV     kvstore.KVStore
	rotationKV kvstore.KVStore
	shiftKV    kvstore.KVStore
	Logger     bot.Logger
}

func NewPluginStore(api plugin.API, logger bot.Logger) Store {
	basicKV := kvstore.NewPluginStore(api)
	return &pluginStore{
		basicKV:    basicKV,
		userKV:     kvstore.NewHashedKeyStore(basicKV, UserKeyPrefix),
		rotationKV: kvstore.NewHashedKeyStore(basicKV, RotationKeyPrefix),
		shiftKV:    kvstore.NewHashedKeyStore(basicKV, ShiftKeyPrefix),
		Logger:     logger,
	}
}
