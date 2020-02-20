// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IndexStore interface {
	Load() (*types.Set, error)
	Store(index *types.Set) error
	DeleteFrom(id string) error
	AddTo(id string) error
}

type indexStore struct {
	key string
	kv  KVStore
}

func (s *store) Index(key string) IndexStore {
	return &indexStore{
		key: key,
		kv:  s.KVStore,
	}
}

func (s *indexStore) Load() (*types.Set, error) {
	index := types.NewSet()
	err := LoadJSON(s.kv, s.key, &index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (s *indexStore) Store(index *types.Set) error {
	err := StoreJSON(s.kv, s.key, index)
	if err != nil {
		return err
	}
	return nil
}

func (s *indexStore) DeleteFrom(id string) error {
	index := types.NewSet()
	err := LoadJSON(s.kv, s.key, &index)
	if err != nil {
		return err
	}

	index.Delete(id)

	err = StoreJSON(s.kv, s.key, index)
	if err != nil {
		return err
	}
	return nil
}

func (s *indexStore) AddTo(id string) error {
	index := types.NewSet()
	err := LoadJSON(s.kv, s.key, &index)
	if err != nil {
		return err
	}

	index.Add(id)

	err = StoreJSON(s.kv, s.key, index)
	if err != nil {
		return err
	}
	return nil
}
