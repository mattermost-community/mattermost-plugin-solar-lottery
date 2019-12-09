// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store/mock_store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

// A rotation requires 3 users, 4 different skills.

var lastRotation0 = map[string]int{"test-rotation": 0}
var lastRotation1 = map[string]int{"test-rotation": 1}
var lastRotation2 = map[string]int{"test-rotation": 2}

var testUserGuru = &store.User{
	MattermostUserID: "test-user-guru",
	SkillLevels: map[string]int{
		"webapp":  4,
		"server":  4,
		"plugins": 4,
	},
	Rotations: lastRotation0,
}

var testUserServer1 = &store.User{
	MattermostUserID: "test-user-server1",
	SkillLevels: map[string]int{
		"webapp":  1,
		"server":  3,
		"plugins": 1,
	},
	Rotations: lastRotation0,
}

var testUserServer2 = &store.User{
	MattermostUserID: "test-user-server2",
	SkillLevels: map[string]int{
		"webapp":  2,
		"server":  3,
		"plugins": 1,
	},
	Rotations: lastRotation0,
}

var testUserServer3 = &store.User{
	MattermostUserID: "test-user-server3",
	SkillLevels: map[string]int{
		"webapp":  1,
		"server":  3,
		"plugins": 1,
	},
	Rotations: lastRotation0,
}

var testUserWebapp1 = &store.User{
	MattermostUserID: "test-user-webapp1",
	SkillLevels: map[string]int{
		"webapp": 3,
		"server": 1,
	},
	Rotations: lastRotation0,
}

var testUserWebapp2 = &store.User{
	MattermostUserID: "test-user-webapp2",
	SkillLevels: map[string]int{
		"webapp": 2,
		"server": 1,
	},
	Rotations: lastRotation0,
}

var testUserWebapp3 = &store.User{
	MattermostUserID: "test-user-webapp3",
	SkillLevels: map[string]int{
		"webapp": 3,
		"server": 1,
	},
	Rotations: lastRotation0,
}

var testUserMobile1 = &store.User{
	MattermostUserID: "test-user-mobile1",
	SkillLevels: map[string]int{
		"webapp": 1,
		"mobile": 3,
	},
	Rotations: lastRotation0,
}

var testUserMobile2 = &store.User{
	MattermostUserID: "test-user-mobile2",
	SkillLevels: map[string]int{
		"webapp": 1,
		"mobile": 3,
	},
	Rotations: lastRotation0,
}

var testUserIDs = store.UserIDList{
	"test-user-guru":    "test-user-guru",
	"test-user-server1": "test-user-server1",
	"test-user-server2": "test-user-server2",
	"test-user-server3": "test-user-server3",
	"test-user-webapp1": "test-user-webapp1",
	"test-user-webapp2": "test-user-webapp2",
	"test-user-webapp3": "test-user-webapp3",
	"test-user-mobile1": "test-user-mobile1",
	"test-user-mobile2": "test-user-mobile2",
}

var testUsers = store.UserList{
	"test-user-guru":    testUserGuru,
	"test-user-server1": testUserServer1,
	"test-user-server2": testUserServer2,
	"test-user-server3": testUserServer3,
	"test-user-webapp1": testUserWebapp1,
	"test-user-webapp2": testUserWebapp2,
	"test-user-webapp3": testUserWebapp3,
	"test-user-mobile1": testUserMobile1,
	"test-user-mobile2": testUserMobile2,
}

var testRotation = store.Rotation{
	Name:              "test-rotation",
	Start:             "2020-01-16",
	Period:            "1m",
	MattermostUserIDs: testUserIDs,
	Size:              3,
	Needs: map[string]store.Need{
		"server-junior": {1, "server", 1},
		"webapp":        {1, "webapp", 2},
		"mobile":        {1, "mobile", 1},
	},
}

func TestPrepareShift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shiftStore := mock_store.NewMockShiftStore(ctrl)
	shiftStore.EXPECT().LoadShift(
		gomock.Eq("test-rotation"),
		gomock.Eq(0),
	).Return(nil, store.ErrNotFound)

	userStore := mock_store.NewMockUserStore(ctrl)
	userStore.EXPECT().LoadUser(
		gomock.Any(),
	).AnyTimes().DoAndReturn(
		func(id string) (*store.User, error) {
			v, ok := testUsers[id]
			if !ok {
				return nil, store.ErrNotFound
			}
			return v, nil
		})

	apiConfig := Config{
		Dependencies: &Dependencies{
			UserStore:  userStore,
			ShiftStore: shiftStore,
			Logger:     &bot.TestLogger{TB: t},
		},
		Config: &config.Config{},
	}

	api := New(apiConfig, "test-user-1").(*api)

	shift, err := api.prepareShift(&testRotation, 0)
	require.Nil(t, err)
	assert.NotNil(t, shift)
}
