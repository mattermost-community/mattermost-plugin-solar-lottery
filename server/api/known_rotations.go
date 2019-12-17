// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

var ErrMultipleResults = errors.New("multiple resolts found")

func withKnownRotations(api *api) error {
	if api.knownRotations != nil {
		return nil
	}

	rr, err := api.RotationStore.LoadKnownRotations()
	if err != nil {
		if err == store.ErrNotFound {
			rr = store.IDMap{}
		} else {
			return err
		}
	}

	api.knownRotations = rr
	return nil
}
