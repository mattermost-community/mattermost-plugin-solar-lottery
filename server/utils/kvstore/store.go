// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"

type Store interface {
	KVStore

	Entity(string) EntityStore
	Index(string, types.IndexPrototype) IndexStore
}

type store struct {
	KVStore
}

func NewStore(kv KVStore) Store {
	return &store{
		KVStore: kv,
	}
}
