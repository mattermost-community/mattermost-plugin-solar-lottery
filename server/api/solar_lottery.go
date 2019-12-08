// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"math"
	"math/rand"
	"sort"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

var ErrShiftAlreadyExists = errors.New("Shift already exists")

func (api *api) scheduleShift(r *store.Rotation, shiftNumber int) (*store.Shift, error) {
	shift, err := api.ShiftStore.LoadShift(r.Name, shiftNumber)
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}

	start, end, err := ShiftDates(r, shiftNumber)
	if shift == nil {
		shift := store.NewShift(r.Name, shiftNumber)
		shift.Start = start.Format(DateFormat)
		shift.End = end.Format(DateFormat)
	} else {
		if shift.Start != start.Format(DateFormat) || shift.End != end.Format(DateFormat) {
			return nil, errors.Errorf("loaded shift has wrong dates %v-%v, expected %v-%v",
				shift.Start, shift.End, start, end)
		}
	}

	rotationUsers, err := api.loadUsers(r.MattermostUserIDs)
	if err != nil {
		return nil, err
	}

	// Filter out users who will not participate
	pool := store.UserList{}
	for _, u := range rotationUsers {
		// Do not include if already in the shift
		_, ok := shift.MattermostUserIDs[u.MattermostUserID]
		if ok {
			continue
		}

		// TODO Do not include if unavailable
		// TODO Do not include if back-to-back check fails

		pool[u.MattermostUserID] = u
	}

	// Try to satisfy the rotation needs
	satisfied := false
	selectedUsers, err := api.loadUsers(shift.MattermostUserIDs)
	for _, need := range r.Needs {
		if len(selectedUsers) >= r.Size {
			break
		}
		need.Count = unsatisfied(need, selectedUsers)
		if need.Count == 0 {
			continue
		}

		picked, err := pickUsers(r, need, pool, shiftNumber)
		if err != nil {
			return nil, err
		}

		for _, u := range picked {
			selectedUsers[u.MattermostUserID] = u
		}
	}

	if !satisfied {
		return nil, errors.New("<><> impossible")
	}
	shift = &store.Shift{}

	return nil, nil
}

func unsatisfied(need store.Need, users store.UserList) int {
	c := need.Count
	for _, u := range users {
		skillLevel, _ := u.SkillLevels[need.Skill]
		if skillLevel < need.Level {
			continue
		}

		c--
		if c == 0 {
			return 0
		}
	}
	return c
}

func pickUsers(r *store.Rotation, need store.Need, users store.UserList, shiftNumber int) (store.UserList, error) {
	skilled := store.UserList{}
	for _, u := range users {
		skillLevel, _ := u.SkillLevels[need.Skill]
		if skillLevel >= need.Level {
			skilled[u.MattermostUserID] = u
		}
	}

	picked := store.UserList{}
	for c := need.Count; c > 0; c-- {
		pickOne(r, need, skilled, shiftNumber)
	}

	return picked, nil
}

func pickOne(r *store.Rotation, need store.Need, users store.UserList, shiftNumber int) (*store.User, error) {
	ids := []string{}
	weights := []float64{}
	for _, user := range users {
		ids = append(ids, user.MattermostUserID)
		weights = append(weights, userWeight(r, user, shiftNumber))
	}

	cdf := make([]float64, len(weights))
	floats.CumSum(cdf, weights)

	random := rand.Float64() * cdf[len(cdf)-1]
	i := sort.Search(len(cdf), func(i int) bool {
		return cdf[i] > random
	})

	return users[ids[i]], nil
}

func userWeight(r *store.Rotation, user *store.User, shiftNumber int) float64 {
	last, _ := user.Joined[r.Name]
	return math.Pow(2.0, float64(shiftNumber-last))
}
