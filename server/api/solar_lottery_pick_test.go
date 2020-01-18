// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPickUser(t *testing.T) {
	const (
		sampleSize = 10000
		low        = 0.9
		high       = 1.1
	)

	for _, tc := range []struct {
		name                string
		users               UserMap
		expectedPercentages map[string]float64
	}{
		{
			name:  "1-way",
			users: usermap(testUserMobile1.withWeight(1e20), testUserMobile2),
			expectedPercentages: map[string]float64{
				testUserMobile1.MattermostUserID: 1,
				testUserMobile2.MattermostUserID: 0,
			},
		},
		{
			name: "fair weighted",
			users: usermap(testUserGuru.withWeight(64),
				testUserServer1.withWeight(32),
				testUserServer2.withWeight(32),
				testUserServer3.withWeight(16),
				testUserMobile1.withWeight(16),
			),
			expectedPercentages: map[string]float64{
				testUserGuru.MattermostUserID:    .4,
				testUserServer1.MattermostUserID: .2,
				testUserServer2.MattermostUserID: .2,
				testUserServer3.MattermostUserID: .1,
				testUserMobile1.MattermostUserID: .1,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counters := map[string]int{}
			for i := 0; i < sampleSize; i++ {
				picked := pickUser(tc.users)
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
