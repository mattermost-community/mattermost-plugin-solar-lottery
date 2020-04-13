// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"sort"
	"testing"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/stretchr/testify/require"
)

func TestSorter(t *testing.T) {
	for _, tc := range []struct {
		name            string
		ids             []types.ID
		weights         []float64
		expectedIDs     []types.ID
		expectedWeights []float64
	}{
		{
			name:            "happy1",
			ids:             []types.ID{"id0", "id1", "id2", "id3", "id4"},
			weights:         []float64{0.01, 3, .3, 7.5, 1e10},
			expectedIDs:     []types.ID{"id4", "id3", "id1", "id2", "id0"},
			expectedWeights: []float64{1e10, 7.5, 3, .3, 0.01},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ws := &sorter{
				Weights: tc.weights,
				IDs:     tc.ids,
			}

			sort.Sort(ws)

			require.Equal(t, tc.expectedIDs, ws.IDs)
			require.Equal(t, tc.expectedWeights, ws.Weights)
		})
	}
}
