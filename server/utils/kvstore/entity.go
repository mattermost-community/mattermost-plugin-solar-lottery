// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"errors"
	"regexp"

	"github.com/mattermost/mattermost-server/v5/model"
)

type EntityStore interface {
	Delete(string) error
	Load(string, interface{}) error
	NewID(name string) (string, error)
	Store(string, interface{}) error
}

type entityStore struct {
	kv KVStore
}

func (s *store) Entity(prefix string) EntityStore {
	return &entityStore{
		kv: NewHashedKeyStore(s.KVStore, prefix+"_"),
	}
}

func (s *entityStore) Load(id string, ref interface{}) error {
	return LoadJSON(s.kv, id, ref)
}

func (s *entityStore) Store(id string, ref interface{}) error {
	return StoreJSON(s.kv, id, ref)
}

func (s *entityStore) Delete(id string) error {
	return s.kv.Delete(id)
}

var ErrTryAgain = errors.New("try again")

func (e *entityStore) NewID(name string) (string, error) {
	for i := 0; i < 5; i++ {
		id := name
		if i > 0 {
			id = name + "-" + model.NewId()[:7]
		}

		dummy := struct{}{}
		err := e.Load(id, &dummy)
		if err == ErrNotFound {
			return id, nil
		}
	}

	return "", ErrTryAgain
}

var reModelID = regexp.MustCompile(`-[a-z0-9]{7}$`)

func NameFromID(id string) string {
	if reModelID.MatchString(id) {
		return id[0 : len(id)-8]
	}
	return id
}
