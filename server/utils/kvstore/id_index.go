// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IDIndexStore interface {
	Load() (*types.StringSet, error)
	Store(*types.StringSet) error
	Delete(string) error
	StoreValue(string) error
}

type idIndexStore struct {
	key string
	kv  KVStore
}

func (s *store) IDIndex(key string) IDIndexStore {
	return &idIndexStore{
		key: key,
		kv:  s.KVStore,
	}
}

func (s *idIndexStore) Load() (*types.StringSet, error) {
	set := types.NewStringSet()
	err := LoadJSON(s.kv, s.key, &set)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *idIndexStore) Store(set *types.StringSet) error {
	err := StoreJSON(s.kv, s.key, set)
	if err != nil {
		return err
	}
	return nil
}

func (s *idIndexStore) Delete(id string) error {
	set, err := s.Load()
	if err != nil {
		return err
	}

	set.Delete(id)

	return s.Store(set)
}

func (s *idIndexStore) StoreValue(v string) error {
	set, err := s.Load()
	if err != nil {
		return err
	}

	set.Set(v)

	return s.Store(set)
}
