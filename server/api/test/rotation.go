// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

// A rotation requires 3 users, 4 different skills.

const RotationName = "test-rotation"
const RotationID = "test-rotation-ID"

func GetTestRotation() *api.Rotation {
	return &api.Rotation{
		Rotation: &store.Rotation{
			MattermostUserIDs: store.IDMap{},
			RotationID:        RotationID,
			Name:              RotationName,
			Period:            api.EveryMonth,
		},
	}
}
