// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/floats"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (f *fill) pickUser(from *sl.Users) *sl.User {
	if from.IsEmpty() {
		return nil
	}

	cdf := make([]float64, from.Len())
	weights := []float64{}
	ids := []types.ID{}
	total := float64(0)
	for _, user := range from.AsArray() {
		ids = append(ids, user.MattermostUserID)
		weight := f.userWeightF(user)
		weights = append(weights, weight)
		total += weight
	}
	floats.CumSum(cdf, weights)
	random := f.rand.Float64() * total
	i := sort.Search(len(cdf), func(i int) bool {
		return cdf[i] >= random
	})
	if i < 0 || i >= len(cdf) {
		return nil
	}

	return from.Get(ids[i])
}

const negligibleWeight = float64(1e-12)
const veryLargeWeight = float64(1e12)

func (f *fill) userWeight(user *sl.User) (w float64) {
	w = f.poolWeights[user.MattermostUserID]
	if w > 0 {
		return w
	}
	defer func() { f.poolWeights[user.MattermostUserID] = w }()

	lastServed := user.GetLastServed(f.r)
	if lastServed > f.time {
		// lastServed in the future can happen for newly added users. Return a
		// non-0 but very low number. This way if there are any other eligible
		// users in the pool this one is not selected. Yet, if all users in the
		// pool are new, one of them is picked.
		return negligibleWeight
	}
	return math.Pow(2, float64(f.time-lastServed)/float64(f.doublingPeriod))
}

func (f *fill) pickRequiredNeed() (done bool, picked sl.Need) {
	s := newSorter(0)
	for _, need := range f.require.AsArray() {
		if need.Count() <= 0 {
			f.require.Delete(need.GetID())
			continue
		}
		id, weight := need.GetID(), f.requiredNeedWeight(need)
		s.Append(id, weight)
	}
	if s.Len() == 0 {
		return true, sl.NewNeed(0, sl.AnySkillLevel)
	}
	sort.Sort(s)

	picked = f.require.Get(s.IDs[0])
	return false, picked
}

func (f *fill) trimRequire() {
	for _, need := range f.require.AsArray() {
		if need.Count() <= 0 {
			f.require.Delete(need.GetID())
		}
	}
}

// Counts each user's weight in once for the need itself, and once for each max
// constraint, then averages the number for the need. The idea is that it
// bubbles up hotter needs, particularly those with restrictions and users "past
// due" so that they get filled first.
func (f *fill) requiredNeedWeight(need sl.Need) float64 {
	var total float64
	users := f.requirePools[need.GetID()].AsArray()
	if len(users) == 0 {
		return veryLargeWeight
	}
	for _, user := range users {
		w := f.userWeight(user)
		total += w
	}

	// Boost needs that have limits
	if f.limit.Contains(need.ID) {
		total = total * 10
	}

	return total / float64(len(users))
}
