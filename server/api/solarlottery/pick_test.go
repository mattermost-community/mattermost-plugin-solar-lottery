// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/test"
)

func TestPickUser(t *testing.T) {
	const (
		sampleSize = 10000
		low        = 0.9
		high       = 1.1
	)

	for _, tc := range []struct {
		name                string
		users               api.UserMap
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
				af.userWeightF = func(u *api.User) float64 {
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
