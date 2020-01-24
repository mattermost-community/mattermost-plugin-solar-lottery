// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/test"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func TestPickUser(t *testing.T) {
	const (
		sampleSize = 10000
		low        = 0.85
		high       = 1.15
	)

	for _, tc := range []struct {
		name                string
		users               sl.UserMap
		weights             map[string]float64
		expectedPercentages map[string]float64
	}{
		{
			name: "1-way",
			users: test.Usermap(
				test.UserMobile1(),
				test.UserMobile2()),
			weights: map[string]float64{
				test.UserIDMobile1: 1e20,
			},
			expectedPercentages: map[string]float64{
				test.UserIDMobile1: 1,
				test.UserIDMobile2: 0,
			},
		},
		{
			name: "fair weighted",
			users: test.Usermap(
				test.UserGuru(),
				test.UserServer1(),
				test.UserServer2(),
				test.UserServer3(),
				test.UserMobile1(),
			),
			weights: map[string]float64{
				test.UserIDGuru:    64,
				test.UserIDServer1: 32,
				test.UserIDServer2: 32,
				test.UserIDServer3: 16,
				test.UserIDMobile1: 16,
			},
			expectedPercentages: map[string]float64{
				test.UserIDGuru:    .4,
				test.UserIDServer1: .2,
				test.UserIDServer2: .2,
				test.UserIDServer3: .1,
				test.UserIDMobile1: .1,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counters := map[string]int{}
			for i := 0; i < sampleSize; i++ {
				af, err := makeTestAutofill(t, 10, nil, nil, nil, 0)
				require.NoError(t, err)
				origWeightF := af.userWeightF
				af.userWeightF = func(u *sl.User) float64 {
					if tc.weights[u.MattermostUserID] != 0 {
						return tc.weights[u.MattermostUserID]
					}
					return origWeightF(u)
				}

				picked := af.pickUser(tc.users)
				require.NotNil(t, picked)
				counters[picked.MattermostUserID]++
			}

			for id := range tc.users {
				p := float64(counters[id]) / float64(sampleSize)
				require.GreaterOrEqual(t, p, tc.expectedPercentages[id]*low, "id %s, percentage %v, expected %v", id, p, tc.expectedPercentages[id])
				require.LessOrEqual(t, p, tc.expectedPercentages[id]*high, "id %s, percentage %v, expected %v", id, p, tc.expectedPercentages[id])
			}
		})
	}
}

func TestPickNeed(t *testing.T) {
	usersServer1 := test.Usermap(
		test.UserGuru(),
		test.UserServer1(),
		test.UserServer2(),
		test.UserServer3(),
		test.UserWebapp1(),
		test.UserWebapp2(),
		test.UserWebapp3(),
	)
	usersServer2 := test.Usermap(
		test.UserGuru(),
		test.UserServer1(),
		test.UserServer2(),
		test.UserServer3(),
	)
	usersWebapp2 := test.Usermap(
		test.UserGuru(),
		test.UserServer2(),
		test.UserWebapp1(),
		test.UserWebapp2(),
		test.UserWebapp3(),
	)

	for _, tc := range []struct {
		name          string
		requiredNeeds store.Needs
		needPools     map[string]sl.UserMap
		expectedNeed  *store.Need
		expectedPool  sl.UserMap
	}{
		{
			name: "happy 2",
			requiredNeeds: store.Needs{
				test.NeedServer_L1_Min3(),
				test.NeedWebapp_L2_Min1(),
			},
			needPools: map[string]sl.UserMap{
				test.NeedServer_L1_Min3().SkillLevel(): usersServer1,
				test.NeedWebapp_L2_Min1().SkillLevel(): usersWebapp2,
			},
			expectedNeed: test.NeedServer_L1_Min3(),
			expectedPool: usersServer1,
		},
		{
			name: "happy 3",
			requiredNeeds: store.Needs{
				test.NeedServer_L1_Min3(),
				test.NeedServer_L2_Min2(),
				test.NeedWebapp_L2_Min1(),
			},
			needPools: map[string]sl.UserMap{
				test.NeedServer_L1_Min3().SkillLevel(): usersServer1,
				test.NeedServer_L2_Min2().SkillLevel(): usersServer2,
				test.NeedWebapp_L2_Min1().SkillLevel(): usersWebapp2,
			},

			// testNeedServer_L2_Min2 is selected since it has the lowest weight/headcount
			expectedNeed: test.NeedServer_L2_Min2(),
			expectedPool: usersServer2,
		},
		{
			name: "happy 3 with constraint",
			requiredNeeds: store.Needs{
				test.NeedServer_L1_Min3(),
				test.NeedServer_L2_Min2(),
				test.NeedWebapp_L2_Min1().WithMax(1),
			},
			needPools: map[string]sl.UserMap{
				test.NeedServer_L1_Min3().SkillLevel(): usersServer1,
				test.NeedServer_L2_Min2().SkillLevel(): usersServer2,
				test.NeedWebapp_L2_Min1().SkillLevel(): usersWebapp2,
			},

			// testNeedServer_L2_Min2 is selected since it has the lowest weight/headcount
			expectedNeed: test.NeedWebapp_L2_Min1().WithMax(1),
			expectedPool: usersWebapp2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			af, err := makeTestAutofill(t, 10, nil, nil, nil, 0)
			require.NoError(t, err)
			need, pool := af.pickNeed(tc.requiredNeeds, tc.needPools)
			require.Equal(t, tc.expectedNeed, need)
			require.Equal(t, tc.expectedPool, pool)
		})
	}
}

func TestUserWeight(t *testing.T) {
	for _, tc := range []struct {
		lastServed     int
		shiftNumber    int
		expectedWeight float64
	}{
		{
			lastServed:     -1,
			shiftNumber:    3,
			expectedWeight: 16.0,
		}, {
			lastServed:     -1,
			shiftNumber:    -1,
			expectedWeight: 1.0,
		}, {
			lastServed:     -1,
			shiftNumber:    -2,
			expectedWeight: 1e-12,
		}, {
			lastServed:     4,
			shiftNumber:    3,
			expectedWeight: 1e-12,
		}, {
			lastServed:     1,
			shiftNumber:    100,
			expectedWeight: 6.338253001141147e+29,
		}, {
			lastServed:     1,
			shiftNumber:    1000,
			expectedWeight: 5.357543035931337e+300,
		}, {
			lastServed:     1,
			shiftNumber:    3,
			expectedWeight: 4.0,
		},
	} {
		t.Run(fmt.Sprintf("%v_%v", tc.lastServed, tc.shiftNumber), func(t *testing.T) {
			af, err := makeTestAutofill(t, 10, nil, nil, nil, tc.shiftNumber)
			require.NoError(t, err)

			user := test.User("test").WithLastServed(test.RotationID, tc.lastServed)

			weight := af.userWeight(user)
			require.Equal(t, tc.expectedWeight, weight)
		})
	}
}
