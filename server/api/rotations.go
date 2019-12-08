// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

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

const Week = time.Hour * 24 * 7
const DateFormat = "2006-01-02"

func ShiftNumber(r *store.Rotation, t time.Time) (int, error) {
	start, err := time.Parse(DateFormat, r.Start)
	if err != nil {
		return 0, err
	}
	if time.Now().Before(start) {
		return 0, errors.Errorf("Time %v is before rotation start %v", t, start)
	}

	switch r.Period {
	case "1w", "w":
		return int(t.Sub(start) / Week), nil
	case "2w":
		return int(t.Sub(start) / (2 * Week)), nil
	case "1m", "m":
		y, m, d := start.Date()
		ty, tm, td := t.Date()
		n := (ty*12 + int(tm)) - (y*12 - int(m))
		if td >= d {
			n++
		}
		return n, nil
	}
	return 0, nil
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
