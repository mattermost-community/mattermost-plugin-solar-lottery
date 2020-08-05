// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"math"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const negligibleWeight = float64(1e-12)
const veryLargeWeight = float64(1e12)

func (f *fill) pickUser(from *sl.Users) *sl.User {
	if from.IsEmpty() {
		return nil
	}

	w := NewWeighted()
	for _, user := range from.AsArray() {
		w.Append(user.MattermostUserID, f.userWeightF(user))
	}

	return from.Get(w.WeightedRandom(f.rand))
}

func (f *fill) pickRequireRandom() (done bool, picked *sl.Need) {
	return f.pickRequireImpl(func(w *weighted) types.ID {
		return w.Random(f.rand)
	})
}

func (f *fill) pickRequireHighest() (done bool, picked *sl.Need) {
	return f.pickRequireImpl(func(w *weighted) types.ID {
		return w.Highest()
	})
}

func (f *fill) pickRequireWeightedRandom() (done bool, picked *sl.Need) {
	return f.pickRequireImpl(func(w *weighted) types.ID {
		return w.WeightedRandom(f.rand)
	})
}

func (f *fill) pickRequireImpl(idf func(*weighted) types.ID) (done bool, picked *sl.Need) {
	w := NewWeighted()
	for _, need := range f.require.AsArray() {
		if need.Count() <= 0 {
			f.require.Delete(need.GetID())
			continue
		}
		id, weight := need.GetID(), f.requireWeight(need)
		if id == sl.AnySkill {
			weight = negligibleWeight
		}
		w.Append(id, weight)
	}
	if w.Len() == 0 {
		return true, nil
	}

	need := f.require.Get(idf(w))
	return false, &need
}

func (f *fill) trimRequire() {
	for _, need := range f.require.AsArray() {
		if need.Count() <= 0 {
			f.require.Delete(need.GetID())
		}
	}
}

func (f *fill) userWeight(user *sl.User) (w float64) {
	w = f.poolWeights[user.MattermostUserID]
	if w > 0 {
		return w
	}
	defer func() { f.poolWeights[user.MattermostUserID] = w }()

	lastServed := user.LastServed.Get(f.r.RotationID)
	if lastServed <= 0 {
		lastServed = f.r.FillSettings.Beginning.Unix()
	}

	if lastServed > f.forTime {
		// lastServed in the future can happen for newly added users. Return a
		// non-0 but very low number. This way if there are any other eligible
		// users in the pool this one is not selected. Yet, if all users in the
		// pool are new, one of them is picked.
		return negligibleWeight
	}
	return math.Pow(2, float64(f.forTime-lastServed)/float64(f.doublingPeriod))
}

// Counts each user's weight in once for the need itself, and once for each max
// constraint, then averages the number for the need. The idea is that it
// bubbles up hotter needs, particularly those with restrictions and users "past
// due" so that they get filled first.
func (f *fill) requireWeight(need sl.Need) float64 {
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
		total *=  10
	}

	return total / float64(len(users))
}
