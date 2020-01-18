// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

// A rotation requires 3 users, 4 different skills.

const testRotationName = "test-rotation"
const testRotationID = "test-rotation-ID"

var testRotation = &Rotation{
	Rotation: &store.Rotation{
		MattermostUserIDs: store.IDMap{},
		RotationID:        testRotationID,
		Name:              testRotationName,
		Period:            EveryMonth,
	},
}

func (rotation *Rotation) withUsers(users UserMap) *Rotation {
	newRotation := rotation.Clone(true)
	newRotation.MattermostUserIDs = make(store.IDMap)
	for id := range users {
		rotation.MattermostUserIDs[id] = id
	}
	newRotation.Users = users
	return newRotation
}
