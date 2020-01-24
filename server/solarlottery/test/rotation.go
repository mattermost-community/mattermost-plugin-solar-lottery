// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

// A rotation requires 3 users, 4 different skills.

const RotationName = "test-rotation"
const RotationID = "test-rotation-ID"

func GetTestRotation() *sl.Rotation {
	return &sl.Rotation{
		Rotation: &store.Rotation{
			MattermostUserIDs: store.IDMap{},
			RotationID:        RotationID,
			Name:              RotationName,
			Period:            sl.EveryMonth,
		},
	}
}
