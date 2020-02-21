// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

type Store interface {
	KVStore

	Entity(string) EntityStore
	Index(string) IndexStore
}

type store struct {
	KVStore
}

func NewStore(kv KVStore) Store {
	return &store{
		KVStore: kv,
	}
}
