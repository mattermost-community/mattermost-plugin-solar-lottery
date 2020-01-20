// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

const Type = "solar-lottery"

type autofiller struct{}

var _ api.Autofiller = (*autofiller)(nil)

func New(logger bot.Logger) api.Autofiller {
	return &autofiller{}
}

// FillShift automatically fills the shift. The caller (api.Guess) is supposed
// to have fully expanded, and deep-cloned the original rotation, so its data is
// not modified. FillShift shallow-clones rotation.Users to preserve the orinal
// map intact, but when called for a sequence of shifts, it relies on the caller
// to carry the users from one call to the next, presumably by using the same
// rotation object.
func (*autofiller) FillShift(rotation *api.Rotation, shiftNumber int, shift *api.Shift, logger bot.Logger) (api.UserMap, error) {
	af, err := makeAutofill(
		rotation.RotationID,
		rotation.Size,
		rotation.Needs.Clone(),
		rotation.Users.Clone(false),
		rotation.ShiftUsers(shift),
		shiftNumber,
		shift.StartTime,
		shift.EndTime,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return af.fill()
}
