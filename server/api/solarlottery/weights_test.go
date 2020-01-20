// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/test"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

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

func TestWeightedUserSorter(t *testing.T) {
	for _, tc := range []struct {
		name            string
		ids             []string
		weights         []float64
		expectedIDs     []string
		expectedWeights []float64
	}{
		{
			name:            "happy1",
			ids:             []string{"id0", "id1", "id2", "id3", "id4"},
			weights:         []float64{0.01, 3, .3, 7.5, 1e10},
			expectedIDs:     []string{"id4", "id3", "id1", "id2", "id0"},
			expectedWeights: []float64{1e10, 7.5, 3, .3, 0.01},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ws := weightedUserSorter{
				weights: tc.weights,
				ids:     tc.ids,
			}

			sort.Sort(&ws)

			require.Equal(t, tc.expectedIDs, ws.ids)
			require.Equal(t, tc.expectedWeights, ws.weights)
		})
	}
}

func TestWeightedNeedSorter(t *testing.T) {
	for _, tc := range []struct {
		name            string
		needs           store.Needs
		weights         []float64
		expectedNeeds   store.Needs
		expectedWeights []float64
	}{
		{
			name: "simple",
			needs: store.Needs{
				&store.Need{Skill: "skill1", Level: 1, Min: 1, Max: -1},
				&store.Need{Skill: "skill2", Level: 2, Min: 2, Max: -1},
				&store.Need{Skill: "skill3", Level: 3, Min: 3, Max: -1},
			},
			weights: []float64{5, 1, 2e10},
			expectedNeeds: store.Needs{
				&store.Need{Skill: "skill3", Level: 3, Min: 3, Max: -1},
				&store.Need{Skill: "skill1", Level: 1, Min: 1, Max: -1},
				&store.Need{Skill: "skill2", Level: 2, Min: 2, Max: -1},
			},
			expectedWeights: []float64{2e10, 5, 1},
		},
		{
			name: "with max 1",
			needs: store.Needs{
				&store.Need{Skill: "skill1", Level: 1, Min: 1, Max: -1},
				&store.Need{Skill: "skill2", Level: 2, Min: 2, Max: 1},
				&store.Need{Skill: "skill3", Level: 3, Min: 3, Max: -1},
			},
			weights: []float64{5, 1, 2e10},
			expectedNeeds: store.Needs{
				&store.Need{Skill: "skill2", Level: 2, Min: 2, Max: 1},
				&store.Need{Skill: "skill3", Level: 3, Min: 3, Max: -1},
				&store.Need{Skill: "skill1", Level: 1, Min: 1, Max: -1},
			},
			expectedWeights: []float64{1, 2e10, 5},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ws := weightedNeedSorter{
				weights: tc.weights,
				needs:   tc.needs,
			}

			sort.Sort(&ws)

			require.Equal(t, tc.expectedNeeds, ws.needs)
			require.Equal(t, tc.expectedWeights, ws.weights)
		})
	}
}
