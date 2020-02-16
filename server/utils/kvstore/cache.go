// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"bytes"
)

type CacheKeyStore struct {
	upstream KVStore

	// key value  of nil indicated deletion
	Data map[string][]byte

	DirtyKeys map[string]bool
	// expires   map[string]time.Time
}

var _ KVStore = (*CacheKeyStore)(nil)

func NewCacheKeyStore(s KVStore) KVStore {
	return &CacheKeyStore{
		upstream: s,
	}
}

func (s *CacheKeyStore) Flush() []error {
	var errs []error
	var err error
	for key := range s.DirtyKeys {
		data := s.Data[key]
		if data == nil {
			err = s.upstream.Delete(key)
		} else {
			err = s.upstream.Store(key, data)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (s *CacheKeyStore) Load(key string) ([]byte, error) {
	data, ok := s.Data[key]
	if ok {
		if data == nil {
			return nil, ErrNotFound
		}
		return data, nil
	}

	data, err := s.upstream.Load(key)
	if err != nil {
		return nil, err
	}

	s.Data[key] = data
	return data, nil
}

func (s *CacheKeyStore) Store(key string, data []byte) error {
	prev, ok := s.Data[key]
	if ok && bytes.Equal(data, prev) {
		return nil
	}

	s.Data[key] = data
	s.DirtyKeys[key] = true
	return nil
}

func (s *CacheKeyStore) StoreTTL(key string, data []byte, ttlSeconds int64) error {
	// TODO Implement expiry
	return s.Store(key, data)
}

func (s *CacheKeyStore) Delete(key string) error {
	return s.Store(key, nil)
}

func (s *CacheKeyStore) Keys() ([]string, error) {
	// Get all keys from the upstream
	keys, err := s.upstream.Keys()
	if err != nil {
		return nil, err
	}

	// Merge with any dirty keys we have
	kmap := map[string]bool{}
	for _, key := range keys {
		kmap[key] = true
	}
	for key := range s.DirtyKeys {
		kmap[key] = true
	}

	// return the merged set
	keys = []string{}
	for key := range kmap {
		keys = append(keys, key)
	}
	return keys, nil
}
