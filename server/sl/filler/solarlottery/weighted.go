// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"math/rand"
	"sort"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"gonum.org/v1/gonum/floats"
)

type weighted struct {
	ids     []types.ID
	weights []float64
	total   float64
}

func NewWeighted() *weighted {
	return &weighted{
		ids:     []types.ID{},
		weights: []float64{},
	}
}

func (w *weighted) Append(id types.ID, weight float64) {
	w.ids = append(w.ids, id)
	w.weights = append(w.weights, weight)
	w.total += weight
}

func (w *weighted) WeightedRandom(rand *rand.Rand) types.ID {
	if len(w.ids) == 0 || len(w.ids) != len(w.weights) {
		return ""
	}

	cdf := make([]float64, len(w.ids))
	floats.CumSum(cdf, w.weights)
	random := rand.Float64() * w.total
	i := sort.Search(len(cdf), func(i int) bool {
		return cdf[i] >= random
	})
	if i < 0 || i >= len(cdf) {
		return ""
	}

	return w.ids[i]
}

func (w *weighted) Random(rand *rand.Rand) types.ID {
	if len(w.ids) == 0 || len(w.ids) != len(w.weights) {
		return ""
	}
	i := rand.Intn(len(w.ids))
	return w.ids[i]
}

func (w *weighted) Highest() types.ID {
	if len(w.ids) == 0 || len(w.ids) != len(w.weights) {
		return ""
	}
	sort.Sort(w)
	return w.ids[0]
}

func (w *weighted) Len() int {
	return len(w.ids)
}

func (w *weighted) Less(i, j int) bool {
	return w.weights[i] > w.weights[j]
}

func (w *weighted) Swap(i, j int) {
	w.weights[i], w.weights[j] = w.weights[j], w.weights[i]
	w.ids[i], w.ids[j] = w.ids[j], w.ids[i]
}
