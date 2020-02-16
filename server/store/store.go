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
	UserKeyPrefix      = "user_"
	RotationKeyPrefix  = "rotation_"
	RecurringKeyPrefix = "recurring_"
	TaskKeyPrefix      = "task_"

	KnownSkillsKey     = "index_skills"
	ActiveRotationsKey = "index_rotations"
)

const OAuth2KeyExpiration = 15 * time.Minute

type Store interface {
	RecurringStore
	RotationStore
	SkillsStore
	TaskStore
	UserStore
}

type pluginStore struct {
	basicKV     kvstore.KVStore
	recurringKV kvstore.KVStore
	rotationKV  kvstore.KVStore
	taskKV      kvstore.KVStore
	userKV      kvstore.KVStore
	Logger      bot.Logger
}

func NewPluginStore(api plugin.API, logger bot.Logger) Store {
	basicKV := kvstore.NewPluginStore(api)
	return &pluginStore{
		basicKV:     basicKV,
		recurringKV: kvstore.NewHashedKeyStore(basicKV, RecurringKeyPrefix),
		rotationKV:  kvstore.NewHashedKeyStore(basicKV, RotationKeyPrefix),
		taskKV:      kvstore.NewHashedKeyStore(basicKV, TaskKeyPrefix),
		userKV:      kvstore.NewHashedKeyStore(basicKV, UserKeyPrefix),
		Logger:      logger,
	}
}
