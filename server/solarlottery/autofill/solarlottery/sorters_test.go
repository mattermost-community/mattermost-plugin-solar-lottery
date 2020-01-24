// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

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
