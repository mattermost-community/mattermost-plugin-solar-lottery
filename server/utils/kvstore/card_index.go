// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type CardIndexStore interface {
	Load() (*types.Index, error)
	Store(*types.Index) error
	Delete(id types.ID) error
	StoreCard(v types.IndexCard) error
}

type cardIndexStore struct {
	key   string
	kv    KVStore
	proto types.IndexCardArray
}

func (s *store) CardIndex(key string, proto types.IndexCardArray) CardIndexStore {
	return &cardIndexStore{
		key:   key,
		kv:    s.KVStore,
		proto: proto,
	}
}

func (s *cardIndexStore) Load() (*types.Index, error) {
	index := types.NewIndex(s.proto)
	err := LoadJSON(s.kv, s.key, &index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (s *cardIndexStore) Store(index *types.Index) error {
	err := StoreJSON(s.kv, s.key, index)
	if err != nil {
		return err
	}
	return nil
}

func (s *cardIndexStore) Delete(id types.ID) error {
	index, err := s.Load()
	if err != nil {
		return err
	}

	index.Delete(id)

	return s.Store(index)
}

func (s *cardIndexStore) StoreCard(v types.IndexCard) error {
	index, err := s.Load()
	if err != nil {
		return err
	}

	index.Set(v)

	return s.Store(index)
}
