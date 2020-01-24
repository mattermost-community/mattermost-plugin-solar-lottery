// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"math"
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/floats"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (af *fill) pickUser(from sl.UserMap) *sl.User {
	if len(from) == 0 {
		return nil
	}

	cdf := make([]float64, len(from))
	weights := []float64{}
	ids := []string{}
	total := float64(0)
	for id, user := range from {
		ids = append(ids, id)
		weight := af.userWeightF(user)
		weights = append(weights, weight)
		total += weight
	}
	floats.CumSum(cdf, weights)
	random := rand.Float64() * total
	i := sort.Search(len(cdf), func(i int) bool {
		return cdf[i] >= random
	})
	if i < 0 || i >= len(cdf) {
		return nil
	}

	return from[ids[i]]
}

func (af *fill) userWeight(user *sl.User) float64 {
	lastServed := user.LastServed[af.rotationID]
	if lastServed > af.shiftNumber {
		// return a non-0 but very low number, to prevent user from serving
		// other than if all have 0 weights
		return 1e-12
	}
	return math.Pow(2.0, float64(af.shiftNumber-lastServed))
}

func (af *fill) pickNeed(requiredNeeds store.Needs, needPools map[string]sl.UserMap) (*store.Need, sl.UserMap) {
	if len(requiredNeeds) == 0 {
		return nil, nil
	}

	s := &weightedNeedSorter{
		weights: make([]float64, len(requiredNeeds)),
		needs:   make(store.Needs, len(requiredNeeds)),
	}
	i := 0
	for _, need := range requiredNeeds {
		pool := needPools[need.SkillLevel()]
		if len(pool) == 0 {
			// not a real possibility, but if there is a need with no pool, declare it the hottest.
			return need, pool
		}
		for _, user := range pool {
			s.weights[i] += af.userWeightF(user)
		}

		// Normalize per remaining needed user, in reverse order, so 1/x.
		s.weights[i] = float64(need.Min) / s.weights[i]
		s.needs[i] = need
		i++
	}

	sort.Sort(s)
	need := s.needs[0]
	pool := needPools[need.SkillLevel()]
	return need, pool
}
