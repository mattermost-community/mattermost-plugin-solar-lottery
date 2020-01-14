// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (api *api) ListShifts(rotation *Rotation, shiftNumber, numShifts int) ([]*Shift, error) {
	err := api.Filter(
	// withActingUser,
	)
	if err != nil {
		return nil, err
	}

	shifts := []*Shift{}
	for i := shiftNumber; i < shiftNumber+numShifts; i++ {
		var shift *Shift
		shift, err = api.loadShift(rotation, i)
		if err != nil {
			if err != store.ErrNotFound {
				return nil, err
			}
			continue
		}
		shifts = append(shifts, shift)
	}

	return shifts, nil
}
