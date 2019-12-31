// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/mock_api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store/mock_store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

// A rotation requires 3 users, 4 different skills.

const testRotation = "test-rotation"
const testRotationID = "test-rotation-ID"

var lastRotation0 = store.IntMap{testRotation: 0}
var lastRotation1 = store.IntMap{testRotation: 1}
var lastRotation2 = store.IntMap{testRotation: 2}

var testUserGuru = makeUser(&store.User{
	MattermostUserID: "test-user-guru",
	SkillLevels: store.IntMap{
		"webapp":  4,
		"server":  4,
		"plugins": 4,
	},
	NextRotationShift: lastRotation0,
})

var testUserServer1 = makeUser(&store.User{
	MattermostUserID: "test-user-server1",
	SkillLevels: store.IntMap{
		"webapp":  1,
		"server":  3,
		"plugins": 1,
	},
	NextRotationShift: lastRotation0,
})

var testUserServer2 = makeUser(&store.User{
	MattermostUserID: "test-user-server2",
	SkillLevels: store.IntMap{
		"webapp":  2,
		"server":  3,
		"plugins": 1,
	},
	NextRotationShift: lastRotation0,
})

var testUserServer3 = makeUser(&store.User{
	MattermostUserID: "test-user-server3",
	SkillLevels: store.IntMap{
		"webapp":  1,
		"server":  3,
		"plugins": 1,
	},
	NextRotationShift: lastRotation0,
})

var testUserWebapp1 = makeUser(&store.User{
	MattermostUserID: "test-user-webapp1",
	SkillLevels: store.IntMap{
		"webapp": 3,
		"server": 1,
	},
	NextRotationShift: lastRotation0,
})

var testUserWebapp2 = makeUser(&store.User{
	MattermostUserID: "test-user-webapp2",
	SkillLevels: store.IntMap{
		"webapp": 2,
		"server": 1,
	},
	NextRotationShift: lastRotation0,
})

var testUserWebapp3 = makeUser(&store.User{
	MattermostUserID: "test-user-webapp3",
	SkillLevels: store.IntMap{
		"webapp": 3,
		"server": 1,
	},
	NextRotationShift: lastRotation0,
})

var testUserMobile1 = makeUser(&store.User{
	MattermostUserID: "test-user-mobile1",
	SkillLevels: store.IntMap{
		"webapp": 1,
		"mobile": 3,
	},
	NextRotationShift: lastRotation0,
})

var testUserMobile2 = makeUser(&store.User{
	MattermostUserID: "test-user-mobile2",
	SkillLevels: store.IntMap{
		"webapp": 1,
		"mobile": 3,
	},
	NextRotationShift: lastRotation0,
})

var testUsers = UserMap{
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

func makeUser(su *store.User) *User {
	return &User{
		User: su,
	}
}

func setupRotationUsers(rotation *Rotation, users ...*User) UserMap {
	rotation.MattermostUserIDs = make(store.IDMap)
	userMap := UserMap{}
	for _, u := range users {
		u = u.Clone()
		rotation.MattermostUserIDs[u.MattermostUserID] = u.MattermostUserID
		userMap[u.MattermostUserID] = u
	}
	return userMap
}

func setupAPIForCrystalBall(t testing.TB, ctrl *gomock.Controller, rotation *Rotation, testUsers UserMap) *api {
	shiftStore := mock_store.NewMockShiftStore(ctrl)
	shiftStore.EXPECT().LoadShift(
		gomock.Eq(rotation.RotationID),
		gomock.Any(),
	).AnyTimes().Return(nil, store.ErrNotFound)

	userStore := mock_store.NewMockUserStore(ctrl)
	userStore.EXPECT().LoadUser(gomock.Any()).AnyTimes().DoAndReturn(
		func(id string) (*store.User, error) {
			user, ok := testUsers[id]
			if !ok {
				return nil, store.ErrNotFound
			}
			return user.User, nil
		})

	rotationStore := mock_store.NewMockRotationStore(ctrl)
	rotationStore.EXPECT().LoadKnownRotations().AnyTimes().Return(store.IDMap{rotation.Name: store.NotEmpty}, nil)
	rotationStore.EXPECT().LoadRotation(rotation.RotationID).AnyTimes().Return(rotation.Rotation, nil)

	pluginAPI := mock_api.NewMockPluginAPI(ctrl)
	pluginAPI.EXPECT().GetMattermostUser(gomock.Any()).AnyTimes().DoAndReturn(
		func(id string) (*model.User, error) {
			user, ok := testUsers[id]
			if !ok {
				return nil, store.ErrNotFound
			}
			return &model.User{
				Id:        user.MattermostUserID,
				Username:  "user" + user.MattermostUserID,
				FirstName: "first-" + user.MattermostUserID,
				LastName:  "last-" + user.MattermostUserID,
			}, nil
		})

	apiConfig := Config{
		Dependencies: &Dependencies{
			UserStore:     userStore,
			ShiftStore:    shiftStore,
			RotationStore: rotationStore,
			PluginAPI:     pluginAPI,
			// Uncomment to display logs while debugging tests
			// Logger: &bot.TestLogger{TB: t},
			Logger: &bot.NilLogger{},
		},
		Config: &config.Config{},
	}

	actingMattermostUserID := "uninitialized"
	for id := range testUsers {
		actingMattermostUserID = id
		break
	}
	api := New(apiConfig, actingMattermostUserID).(*api)
	return api
}

func need(c int, skill string, level int) store.Need {
	return store.Need{
		Min:   c,
		Skill: skill,
		Level: level,
	}
}

func needOne(skill string, level int) store.Need {
	return need(1, skill, level)
}

func TestPrepareShiftHappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rotation := &Rotation{
		Rotation: &store.Rotation{
			RotationID: testRotationID,
			Name:       testRotation,
			Start:      "2020-01-16",
			Period:     EveryMonth,
			Size:       3,
			Needs: []store.Need{
				needOne("server", 1),
				needOne("webapp", 2),
				needOne("mobile", 1),
			},
		},
		Users: UserMap{},
	}

	testUsers := setupRotationUsers(rotation,
		testUserGuru, testUserServer1, testUserServer2, testUserServer3, testUserWebapp1,
		testUserWebapp2, testUserWebapp3, testUserMobile1, testUserMobile2)

	api := setupAPIForCrystalBall(t, ctrl, rotation, testUsers)

	shifts, err := api.Guess(rotation, 0, 1, true)
	require.Nil(t, err)
	assert.Len(t, shifts, 1)
	require.Equal(t, rotation.Size, len(shifts[0].MattermostUserIDs))
}

func TestPrepareShiftEvenDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rotation := &Rotation{
		Rotation: &store.Rotation{
			Name:   testRotation,
			Start:  "2020-01-16",
			Period: EveryMonth,
			Size:   1,
			Needs: []store.Need{
				needOne("webapp", 1),
			},
		},
	}

	testUsers := setupRotationUsers(rotation,
		testUserGuru, testUserServer1, testUserServer2, testUserServer3, testUserWebapp1,
		testUserWebapp2, testUserWebapp3, testUserMobile1, testUserMobile2)

	api := setupAPIForCrystalBall(t, ctrl, rotation, testUsers)

	counters := store.IntMap{}
	shifts, err := api.Guess(rotation, 3, len(testUsers)*1000, true)
	require.Nil(t, err)
	require.Len(t, shifts, len(testUsers)*1000)

	for _, shift := range shifts {
		assert.NotNil(t, shift)
		require.Equal(t, rotation.Size, len(shift.MattermostUserIDs))
		for mattermostUserID := range shift.MattermostUserIDs {
			counters[mattermostUserID]++
		}
	}

	for k, c := range counters {
		assert.Greater(t, c, 900, k)
		assert.Less(t, c, 1100, k)
	}
}
