// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type sorter struct {
	IDs     []types.ID
	Weights []float64
}

func newSorter(size int) *sorter {
	return &sorter{
		IDs:     make([]types.ID, size),
		Weights: make([]float64, size),
	}
}

func (s *sorter) Append(id types.ID, weight float64) {
	s.IDs = append(s.IDs, id)
	s.Weights = append(s.Weights, weight)
}

func (s *sorter) Len() int {
	return len(s.IDs)
}

func (s *sorter) Less(i, j int) bool {
	return s.Weights[i] > s.Weights[j]
}

func (s *sorter) Swap(i, j int) {
	s.Weights[i], s.Weights[j] = s.Weights[j], s.Weights[i]
	s.IDs[i], s.IDs[j] = s.IDs[j], s.IDs[i]
}
