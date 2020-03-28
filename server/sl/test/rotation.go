// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
)

// A rotation requires 3 users, 4 different skills.

const RotationName = "test-rotation"
const RotationID = "test-rotation-ID"

func GetTestRotation() *sl.Rotation {
	r := sl.NewRotation()
	r.RotationID = RotationID
	return r
}
