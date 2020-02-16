// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/autofill/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/mock_solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store/mock_store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func solarLotteryForGuess(t testing.TB, ctrl *gomock.Controller, rotation *sl.Rotation, usersDataSource sl.UserMap) sl.SolarLottery {
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
	rotationStore.EXPECT().LoadActiveRotations().AnyTimes().Return(store.IDMap{rotation.Name: store.NotEmpty}, nil)
	rotationStore.EXPECT().LoadRotation(rotation.RotationID).AnyTimes().Return(rotation.Rotation, nil)

	pluginAPI := mock_solarlottery.NewMockPluginAPI(ctrl)
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

	apiConfig := sl.Config{
		Dependencies: &sl.Dependencies{
			Autofillers: map[string]sl.Autofiller{
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

	return sl.New(apiConfig, actingMattermostUserID)
}

func TestPrepareShiftHappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rotation := GetTestRotation()
	rotation.Period = sl.EveryMonth
	rotation.Size = 3
	rotation.Needs = store.Needs{
		NeedServer_L1_Min1(),
		NeedWebapp_L2_Min1(),
		NeedMobile_L1_Min1(),
	}
	rotation = rotation.WithUsers(AllUsers())
	rotation = rotation.WithStart("2020-01-16")

	sl := solarLotteryForGuess(t, ctrl, rotation, AllUsers())
	shifts, err := sl.Guess(rotation, 0, 1)

	require.Nil(t, err)
	assert.Len(t, shifts, 1)
	require.Equal(t, rotation.Size, len(shifts[0].MattermostUserIDs))
}

func TestPrepareShiftEvenDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rotation := GetTestRotation()
	rotation.Period = sl.EveryMonth
	rotation.Size = 1
	rotation.Needs = store.Needs{
		NeedWebapp_L1_Min1(),
	}
	rotation = rotation.WithUsers(AllUsers())
	rotation = rotation.WithStart("2020-01-16")

	sl := solarLotteryForGuess(t, ctrl, rotation, AllUsers())

	sampleSize := 200
	counters := store.IntMap{}
	shifts, err := sl.Guess(rotation, 3, len(AllUsers())*sampleSize)
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
