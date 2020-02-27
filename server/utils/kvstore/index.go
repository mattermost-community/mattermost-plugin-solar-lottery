// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IndexStore interface {
	Load() (types.Index, error)
	Store(types.Index) error
	Delete(id string) error
	StoreValue(v types.Identifiable) error
}

type indexStore struct {
	key   string
	kv    KVStore
	proto types.IndexPrototype
}

func (s *store) Index(key string, proto types.IndexPrototype) IndexStore {
	return &indexStore{
		key:   key,
		kv:    s.KVStore,
		proto: proto,
	}
}

func (s *indexStore) Load() (types.Index, error) {
	index := types.NewIndex(s.proto)
	err := LoadJSON(s.kv, s.key, &index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (s *indexStore) Store(index types.Index) error {
	err := StoreJSON(s.kv, s.key, index)
	if err != nil {
		return err
	}
	return nil
}

func (s *indexStore) Delete(id string) error {
	index, err := s.Load()
	if err != nil {
		return err
	}

	index.Delete(id)

	return s.Store(index)
}

func (s *indexStore) StoreValue(v types.Identifiable) error {
	index, err := s.Load()
	if err != nil {
		return err
	}

	index.Set(v)

	return s.Store(index)
}
