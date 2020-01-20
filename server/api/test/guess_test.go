// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/mock_api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store/mock_store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func setupAPIForGuess(t testing.TB, ctrl *gomock.Controller, rotation *api.Rotation, usersDataSource api.UserMap) api.API {
	shiftStore := mock_store.NewMockShiftStore(ctrl)
	shiftStore.EXPECT().LoadShift(
		gomock.Eq(rotation.RotationID),
		gomock.Any(),
	).AnyTimes().Return(nil, store.ErrNotFound)

	userStore := mock_store.NewMockUserStore(ctrl)
	userStore.EXPECT().LoadUser(gomock.Any()).AnyTimes().DoAndReturn(
		func(id string) (*store.User, error) {
			user, ok := usersDataSource[id]
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
			user, ok := usersDataSource[id]
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

	// Uncomment to display logs while debugging tests
	// logger := &bot.TestLogger{TB: t}
	logger := &bot.NilLogger{}

	apiConfig := api.Config{
		Dependencies: &api.Dependencies{
			Autofillers: map[string]api.Autofiller{
				"":                solarlottery.New(logger), // default
				solarlottery.Type: solarlottery.New(logger),
			},
			UserStore:     userStore,
			ShiftStore:    shiftStore,
			RotationStore: rotationStore,
			PluginAPI:     pluginAPI,
			Logger:        logger,
		},
		Config: &config.Config{},
	}

	actingMattermostUserID := "uninitialized"
	for id := range usersDataSource {
		actingMattermostUserID = id
		break
	}

	return api.New(apiConfig, actingMattermostUserID)
}

func TestPrepareShiftHappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rotation := GetTestRotation()
	rotation.Period = api.EveryMonth
	rotation.Size = 3
	rotation.Needs = store.Needs{
		NeedServer_L1_Min1(),
		NeedWebapp_L2_Min1(),
		NeedMobile_L1_Min1(),
	}
	rotation = rotation.WithUsers(AllUsers())
	rotation = rotation.WithStart("2020-01-16")

	api := setupAPIForGuess(t, ctrl, rotation, AllUsers())
	shifts, err := api.Guess(rotation, 0, 1)

	require.Nil(t, err)
	assert.Len(t, shifts, 1)
	require.Equal(t, rotation.Size, len(shifts[0].MattermostUserIDs))
}

func TestPrepareShiftEvenDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rotation := GetTestRotation()
	rotation.Period = api.EveryMonth
	rotation.Size = 1
	rotation.Needs = store.Needs{
		NeedWebapp_L1_Min1(),
	}
	rotation = rotation.WithUsers(AllUsers())
	rotation = rotation.WithStart("2020-01-16")

	api := setupAPIForGuess(t, ctrl, rotation, AllUsers())

	sampleSize := 200
	counters := store.IntMap{}
	shifts, err := api.Guess(rotation, 3, len(AllUsers())*sampleSize)
	require.Nil(t, err)
	require.Len(t, shifts, len(AllUsers())*sampleSize)

	for _, shift := range shifts {
		assert.NotNil(t, shift)
		require.Equal(t, rotation.Size, len(shift.MattermostUserIDs))
		for mattermostUserID := range shift.MattermostUserIDs {
			counters[mattermostUserID]++
		}
	}

	for k, c := range counters {
		assert.Greater(t, c, sampleSize*90/100, k)
		assert.Less(t, c, sampleSize*110/100, k)
	}
}
