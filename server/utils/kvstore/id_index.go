// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IDIndexStore interface {
	Load() (*types.IDIndex, error)
	Store(*types.IDIndex) error
	Delete(types.ID) error
	Set(types.ID) error
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

func (s *idIndexStore) Load() (*types.IDIndex, error) {
	set := types.NewIDIndex()
	err := LoadJSON(s.kv, s.key, &set)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *idIndexStore) Store(index *types.IDIndex) error {
	err := StoreJSON(s.kv, s.key, index)
	if err != nil {
		return err
	}
	return nil
}

func (s *idIndexStore) Delete(id types.ID) error {
	index, err := s.Load()
	if err != nil {
		return err
	}

	index.Delete(id)
	return s.Store(index)
}

func (s *idIndexStore) Set(v types.ID) error {
	index, err := s.Load()
	switch err {
	case nil:

	case ErrNotFound:
		index = types.NewIDIndex()

	default:
		return err
	}

	index.Set(v)

	return s.Store(index)
}
