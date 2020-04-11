// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/test"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestPickUser(t *testing.T) {
	const (
		sampleSize = 10000
		low        = 0.85
		high       = 1.15
	)

	for _, tc := range []struct {
		name                string
		users               *sl.Users
		weights             map[types.ID]float64
		expectedPercentages map[types.ID]float64
	}{
		{
			name: "1-way",
			users: sl.NewUsers(
				test.UserMobile1(),
				test.UserMobile2()),
			weights: map[types.ID]float64{
				test.UserIDMobile1: 1e20,
			},
			expectedPercentages: map[types.ID]float64{
				test.UserIDMobile1: 1,
				test.UserIDMobile2: 0,
			},
		},
		{
			name: "fair weighted",
			users: sl.NewUsers(
				test.UserGuru(),
				test.UserServer1(),
				test.UserServer2(),
				test.UserServer3(),
				test.UserMobile1(),
			),
			weights: map[types.ID]float64{
				test.UserIDGuru:    64,
				test.UserIDServer1: 32,
				test.UserIDServer2: 32,
				test.UserIDServer3: 16,
				test.UserIDMobile1: 16,
			},
			expectedPercentages: map[types.ID]float64{
				test.UserIDGuru:    .4,
				test.UserIDServer1: .2,
				test.UserIDServer2: .2,
				test.UserIDServer3: .1,
				test.UserIDMobile1: .1,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counters := map[types.ID]int{}
			for i := 0; i < sampleSize; i++ {
				f := makeTestFiller(t, tc.users, nil, nil, nil)
				origWeightF := f.userWeightF
				f.userWeightF = func(u *sl.User) float64 {
					if tc.weights[u.MattermostUserID] != 0 {
						return tc.weights[u.MattermostUserID]
					}
					return origWeightF(u)
				}

				picked := f.pickUser(tc.users)
				require.NotNil(t, picked)
				counters[picked.MattermostUserID]++
			}

			for _, id := range tc.users.IDs() {
				p := float64(counters[id]) / float64(sampleSize)
				require.GreaterOrEqual(t, p, tc.expectedPercentages[id]*low, "id %s, percentage %v, expected %v", id, p, tc.expectedPercentages[id])
				require.LessOrEqual(t, p, tc.expectedPercentages[id]*high, "id %s, percentage %v, expected %v", id, p, tc.expectedPercentages[id])
			}
		})
	}
}

func TestPickNeed(t *testing.T) {
	usersServer1 := sl.NewUsers(
		test.UserGuru().WithLastServed(test.RotationID, types.MustParseTime("2019-03-01")),
		test.UserServer1(),
		test.UserServer2(),
		test.UserServer3(),
		test.UserWebapp1(),
		test.UserWebapp3(),
	)
	usersServer2 := sl.NewUsers(
		test.UserGuru().WithLastServed(test.RotationID, types.MustParseTime("2019-03-01")),
		test.UserServer1(),
		test.UserServer2(),
		test.UserServer3(),
	)
	usersWebapp2 := sl.NewUsers(
		test.UserGuru().WithLastServed(test.RotationID, types.MustParseTime("2019-03-01")),
		test.UserServer2(),
		test.UserWebapp1(),
		test.UserWebapp2().WithLastServed(test.RotationID, types.MustParseTime("2020-02-01")),
		test.UserWebapp3(),
	)

	for _, tc := range []struct {
		name         string
		require      *sl.Needs
		requirePools map[types.ID]*sl.Users // by need ID (SkillLevel as string)
		limit        *sl.Needs
		time         types.Time
		expectedNeed sl.Need
		expectedDone bool
	}{
		{
			name: "happy 2",
			require: sl.NewNeeds(
				test.C3_Server_L1(),
				test.C1_Webapp_L2(),
			),
			requirePools: map[types.ID]*sl.Users{
				test.C3_Server_L1().GetID(): usersServer1,
				test.C1_Webapp_L2().GetID(): usersWebapp2,
			},
			expectedNeed: test.C1_Webapp_L2(),
		},
		{
			name: "happy 3",
			require: sl.NewNeeds(
				test.C3_Server_L1(),
				test.C2_Server_L2(),
				test.C1_Webapp_L2(),
			),
			requirePools: map[types.ID]*sl.Users{
				test.C3_Server_L1().GetID(): usersServer1,
				test.C2_Server_L2().GetID(): usersServer2,
				test.C1_Webapp_L2().GetID(): usersWebapp2,
			},

			// testNeedServer_L2_Min2 is selected since it has the lowest weight/headcount
			expectedNeed: test.C2_Server_L2(),
		},
		{
			name: "happy 3 with constraint",
			require: sl.NewNeeds(
				test.C3_Server_L1(),
				test.C2_Server_L2(),
				test.C1_Webapp_L2(),
			),
			limit: sl.NewNeeds(
				test.C1_Webapp_L2(),
			),
			requirePools: map[types.ID]*sl.Users{
				test.C3_Server_L1().GetID(): usersServer1,
				test.C2_Server_L2().GetID(): usersServer2,
				test.C1_Webapp_L2().GetID(): usersWebapp2,
			},

			expectedNeed: test.C1_Webapp_L2(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := makeTestFiller(t, nil, nil, tc.require, tc.limit)
			f.requirePools = tc.requirePools
			done, need := f.pickRequiredNeed()
			require.Equal(t, tc.expectedNeed, need)
			require.Equal(t, tc.expectedDone, done)
		})
	}
}

func TestNeedWeight(t *testing.T) {
	userGuru := test.UserGuru().WithLastServed(test.RotationID, types.MustParseTime("2019-03-01"))
	usersServer := sl.NewUsers(userGuru, test.UserServer1(), test.UserServer2())
	usersWebapp := sl.NewUsers(userGuru, test.UserServer2(), test.UserWebapp1(), test.UserWebapp3())

	for _, tc := range []struct {
		name         string
		require      *sl.Needs
		requirePools map[types.ID]*sl.Users // by need ID (SkillLevel as string)
		limit        *sl.Needs
	}{
		{
			name: "happy",
			require: sl.NewNeeds(
				test.C3_Server_L1(),
				test.C1_Webapp_L2(),
			),
			requirePools: map[types.ID]*sl.Users{
				test.C3_Server_L1().GetID(): usersServer,
				test.C1_Webapp_L2().GetID(): usersWebapp,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := makeTestFiller(t, nil, nil, tc.require, tc.limit)
			f.requirePools = tc.requirePools
			w := f.requiredNeedWeight(test.C1_Webapp_L2())
			require.Equal(t, 1.3724873597089898e+15, w)
			w = f.requiredNeedWeight(test.C3_Server_L1())
			require.Equal(t, 1.8299831462785262e+15, w)
		})
	}
}
func TestUserWeight(t *testing.T) {
	for _, tc := range []struct {
		lastServed     types.Time
		time           types.Time
		doubling       int64
		expectedWeight float64
	}{
		{
			lastServed:     types.MustParseTime("2020-03-01"),
			time:           types.MustParseTime("2020-03-18"),
			expectedWeight: 5.383600770529424,
		}, {
			lastServed:     types.MustParseTime("2020-03-01"),
			time:           types.MustParseTime("2020-03-01"),
			expectedWeight: 1,
		}, {
			lastServed:     types.MustParseTime("2020-03-25"),
			time:           types.MustParseTime("2020-03-01"),
			expectedWeight: negligibleWeight,
		}, {
			lastServed:     types.MustParseTime("2020-03-01"),
			time:           types.MustParseTime("2035-03-01"),
			doubling:       28 * 24 * 3600,
			expectedWeight: 7.840945539427057e+58,
		},
	} {
		t.Run(fmt.Sprintf("%v_%v", tc.lastServed, tc.time), func(t *testing.T) {
			f := makeTestFiller(t, nil, nil, nil, nil)
			f.time = tc.time.Unix()
			if tc.doubling > 0 {
				f.doublingPeriod = tc.doubling
			}

			user := sl.NewUser("test").WithLastServed(test.RotationID, tc.lastServed)
			weight := f.userWeight(user)
			require.Equal(t, tc.expectedWeight, weight)
		})
	}
}
