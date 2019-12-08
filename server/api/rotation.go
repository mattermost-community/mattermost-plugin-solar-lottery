// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Rotations interface {
	AddRotation(r *store.Rotation) (*store.Rotation, error)
	ChangeRotationNeed(r *store.Rotation, name, skill string, level, count int)
	DeleteRotation(string) error
	ListRotations() (map[string]*store.Rotation, error)
	RemoveRotationNeed(r *store.Rotation, name string) error
	UpdateRotation(string, func(*store.Rotation) error) (*store.Rotation, error)
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
	rr := api.rotations
	_, ok := rr[r.Name]
	if ok {
		return nil, ErrRotationAlreadyExists
	}

	rr[r.Name] = r
	err = api.RotationsStore.StoreRotations(rr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (api *api) DeleteRotation(rotationName string) error {
	rr, err := api.RotationsStore.LoadRotations()
	if err != nil && err != store.ErrNotFound {
		return err
	}
	_, ok := rr[rotationName]
	if !ok {
		return store.ErrNotFound
	}

	delete(rr, rotationName)
	return api.RotationsStore.StoreRotations(rr)
}

func withRotations(api *api) error {
	if api.rotations != nil {
		return nil
	}

	rr, err := api.RotationsStore.LoadRotations()
	if err != nil {
		if err == store.ErrNotFound {
			rr = map[string]*store.Rotation{}
		} else {
			return err
		}
	}

	api.rotations = rr
	return nil
}
