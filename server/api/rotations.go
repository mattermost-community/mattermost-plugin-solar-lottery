// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Rotations interface {
	ListRotations() (map[string]*store.Rotation, error)
	UpdateRotation(string, func(*store.Rotation) error) (*store.Rotation, error)
	DeleteRotation(string) error
	AddRotation(r *store.Rotation) (*store.Rotation, error)
}

var ErrRotationAlreadyExists = errors.New("rotation already exists")

func (api *api) ListRotations() (map[string]*store.Rotation, error) {
	err := api.Filter(withRotations)
	if err != nil {
		return nil, err
	}
	return api.rotations, nil
}

func (api *api) UpdateRotation(rotationName string, updatef func(*store.Rotation) error) (*store.Rotation, error) {
	err := api.Filter(withRotations)
	if err != nil {
		return nil, err
	}
	r, ok := api.rotations[rotationName]
	if !ok {
		return nil, store.ErrNotFound
	}

	err = updatef(r)
	if err != nil {
		return nil, err
	}
	return r, api.RotationsStore.StoreRotations(api.rotations)
}

func (api *api) AddRotation(r *store.Rotation) (*store.Rotation, error) {
	err := api.Filter(withRotations)
	if err != nil {
		return nil, err
	}
	rotations := api.rotations
	_, ok := rotations[r.Name]
	if ok {
		return nil, ErrRotationAlreadyExists
	}

	rotations[r.Name] = r
	err = api.RotationsStore.StoreRotations(rotations)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (api *api) DeleteRotation(rotationName string) error {
	rotations, err := api.RotationsStore.LoadRotations()
	if err != nil && err != store.ErrNotFound {
		return err
	}
	_, ok := rotations[rotationName]
	if !ok {
		return store.ErrNotFound
	}

	delete(rotations, rotationName)
	return api.RotationsStore.StoreRotations(rotations)
}

func withRotations(api *api) error {
	if api.rotations != nil {
		return nil
	}

	rotations, err := api.RotationsStore.LoadRotations()
	if err != nil {
		if err == store.ErrNotFound {
			rotations = map[string]*store.Rotation{}
		} else {
			return err
		}
	}

	api.rotations = rotations
	return nil
}
