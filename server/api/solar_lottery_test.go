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

func setupRotationUsers(r *store.Rotation, users ...*store.User) store.UserList {
	r.MattermostUserIDs = make(store.UserIDList)
	userList := make(store.UserList)
	for _, u := range users {
		u = u.Clone()
		r.MattermostUserIDs[u.MattermostUserID] = u.MattermostUserID
		userList[u.MattermostUserID] = u
	}
	return userList
}

func setupAPIForPrepareShift(t testing.TB, ctrl *gomock.Controller, testUsers store.UserList) *api {
	shiftStore := mock_store.NewMockShiftStore(ctrl)
	shiftStore.EXPECT().LoadShift(
		gomock.Eq("test-rotation"),
		gomock.Any(),
	).AnyTimes().Return(nil, store.ErrNotFound)

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
			// Uncomment to display logs while debugging tests
			// Logger:     &bot.TestLogger{TB: t},
			Logger: &bot.NilLogger{},
		},
		Config: &config.Config{},
	}

	api := New(apiConfig, "test-user-1").(*api)
	return api
}

func TestPrepareShiftHappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var testRotation = store.Rotation{
		Name:   "test-rotation",
		Start:  "2020-01-16",
		Period: "1m",
		Size:   3,
		Needs: map[string]store.Need{
			"server-junior": {1, "server", 1},
			"webapp":        {1, "webapp", 2},
			"mobile":        {1, "mobile", 1},
		},
	}

	testUsers := setupRotationUsers(&testRotation,
		testUserGuru, testUserServer1, testUserServer2, testUserServer3, testUserWebapp1,
		testUserWebapp2, testUserWebapp3, testUserMobile1, testUserMobile2)

	api := setupAPIForPrepareShift(t, ctrl, testUsers)

	shift, err := api.prepareShift(&testRotation, 0)
	require.Nil(t, err)
	assert.NotNil(t, shift)
	require.Equal(t, testRotation.Size, len(shift.MattermostUserIDs))
}

func TestPrepareShiftEvenDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var testRotation = store.Rotation{
		Name:   "test-rotation",
		Start:  "2020-01-16",
		Period: "1m",
		Size:   1,
		Needs: map[string]store.Need{
			"webapp": {1, "webapp", 1},
		},
	}

	testUsers := setupRotationUsers(&testRotation,
		testUserGuru, testUserServer1, testUserServer2, testUserServer3, testUserWebapp1,
		testUserWebapp2, testUserWebapp3, testUserMobile1, testUserMobile2)

	api := setupAPIForPrepareShift(t, ctrl, testUsers)

	counters := map[string]int{}
	for i := 0; i < len(testUsers)*1000; i++ {
		shift, err := api.prepareShift(&testRotation, i)
		require.Nil(t, err)
		assert.NotNil(t, shift)
		require.Equal(t, testRotation.Size, len(shift.MattermostUserIDs))

		for mattermostUserID := range shift.MattermostUserIDs {
			testUsers[mattermostUserID].Rotations["test-rotation"] = i
			counters[mattermostUserID]++
		}
	}

	for k, c := range counters {
		assert.Greater(t, c, 900, k)
		assert.Less(t, c, 1100, k)
	}
}
