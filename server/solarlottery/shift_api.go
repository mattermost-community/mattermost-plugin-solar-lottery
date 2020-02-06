// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (sl *solarLottery) ListShifts(rotation *Rotation, shiftNumber, numShifts int) ([]*Shift, error) {
	shifts := []*Shift{}
	for i := shiftNumber; i < shiftNumber+numShifts; i++ {
		var shift *Shift
		shift, err := sl.loadShift(rotation, i)
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
