// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type weightedNeedSorter struct {
	needs   store.Needs
	weights []float64
}

func (s *weightedNeedSorter) Len() int {
	return len(s.needs)
}

// Sort all needs that have a max limit to the top, to reduce hitting that
// limit opportunistically. Otherwise sort by normalized weight representing
// staffing heat level.
func (s *weightedNeedSorter) Less(i, j int) bool {
	if (s.needs[i].Max < 0) != (s.needs[j].Max < 0) {
		return s.needs[j].Max < 0
	}
	return s.weights[i] > s.weights[j]
}

func (s *weightedNeedSorter) Swap(i, j int) {
	s.weights[i], s.weights[j] = s.weights[j], s.weights[i]
	s.needs[i], s.needs[j] = s.needs[j], s.needs[i]
}

type weightedUserSorter struct {
	ids     []string
	weights []float64
}

func (s *weightedUserSorter) Len() int {
	return len(s.ids)
}

func (s *weightedUserSorter) Less(i, j int) bool {
	return s.weights[i] > s.weights[j]
}

func (s *weightedUserSorter) Swap(i, j int) {
	s.weights[i], s.weights[j] = s.weights[j], s.weights[i]
	s.ids[i], s.ids[j] = s.ids[j], s.ids[i]
}
