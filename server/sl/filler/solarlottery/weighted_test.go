// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestWeighted(t *testing.T) {
	ids := []types.ID{"0", "1", "2", "3", "4"}
	weights := []float64{.11, 40, 3, 4, 5}
	expectedSorted := []types.ID{"1", "4", "3", "2", "0"}
	expectedHighest := types.ID("1")

	for _, tc := range []struct {
		seq                    int64
		expectedRandom         []types.ID
		expectedWeightedRandom []types.ID
	}{
		{
			seq:                    1,
			expectedRandom:         []types.ID{"1", "2", "1", "0", "1", "4", "2", "3", "1", "2"},
			expectedWeightedRandom: []types.ID{"4", "1", "1", "1", "1", "2", "1", "1", "1", "1"},
		},
		{
			seq:                    2,
			expectedRandom:         []types.ID{"1", "2", "4", "0", "1", "0", "0", "1", "0", "4"},
			expectedWeightedRandom: []types.ID{"1", "1", "4", "1", "1", "3", "1", "4", "1", "3"},
		},
		{
			seq:                    3,
			expectedRandom:         []types.ID{"3", "1", "2", "1", "4", "1", "1", "1", "4", "4"},
			expectedWeightedRandom: []types.ID{"1", "1", "1", "1", "1", "2", "1", "1", "1", "1"},
		},
		{
			seq:                    4,
			expectedRandom:         []types.ID{"4", "3", "2", "4", "1", "2", "0", "0", "4", "0"},
			expectedWeightedRandom: []types.ID{"1", "4", "1", "1", "1", "1", "1", "3", "1", "1"},
		},
	} {
		t.Run(strconv.FormatInt(tc.seq, 10), func(t *testing.T) {
			rr := rand.New(rand.NewSource(tc.seq))
			w := NewWeighted()
			for i, id := range ids {
				w.Append(id, weights[i])
			}

			random := []types.ID{}
			weightedRandom := []types.ID{}
			for i := 0; i < 10; i++ {
				random = append(random, w.Random(rr))
				weightedRandom = append(weightedRandom, w.WeightedRandom(rr))
			}
			require.Equal(t, tc.expectedRandom, random)
			require.Equal(t, tc.expectedWeightedRandom, weightedRandom)
			sort.Sort(w)
			require.Equal(t, expectedSorted, w.ids)
			require.Equal(t, expectedHighest, w.Highest())
		})
	}
}
