// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

const Type = "solar-lottery"

type autofiller struct{}

var _ sl.Autofiller = (*autofiller)(nil)

func New(logger bot.Logger) sl.Autofiller {
	return &autofiller{}
}

// FillShift automatically fills the shift. The caller (sl.Guess) is supposed
// to have fully expanded, and deep-cloned the original rotation, so its data is
// not modified. FillShift shallow-clones rotation.Users to preserve the orinal
// map intact, but when called for a sequence of shifts, it relies on the caller
// to carry the users from one call to the next, presumably by using the same
// rotation object.
func (*autofiller) FillShift(rotation *sl.Rotation, shiftNumber int, shift *sl.Shift, logger bot.Logger) (sl.UserMap, error) {
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
