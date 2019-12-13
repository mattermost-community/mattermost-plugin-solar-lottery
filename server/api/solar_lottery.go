// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type SolarLottery interface {
	CrystalBall(rotationName string, startShift int, numShifts int, autofill bool) ([]*store.Shift, error)
}

var ErrShiftAlreadyExists = errors.New("Shift already exists")

func (api *api) CrystalBall(rotationName string, startShift int, numShifts int, autofill bool) ([]*store.Shift, error) {
	err := api.Filter(withRotations)
	if err != nil {
		return nil, err
	}

	r := api.rotations[rotationName]
	if r == nil {
		return nil, store.ErrNotFound
	}

	// Load all rotation users
	rotationUsers, err := api.loadUsers(r.MattermostUserIDs)
	if err != nil {
		return nil, err
	}

	cachedUsers := store.UserList{}
	for k, v := range rotationUsers {
		cachedUsers[k] = v
	}

	var shifts []*store.Shift
	for n := startShift; n < numShifts; n++ {
		start, end, err := ShiftDates(r, n)
		if err != nil {
			return nil, err
		}

		var shift *store.Shift
		shift, err = api.ShiftStore.LoadShift(r.Name, n)
		switch {
		case autofill && err == store.ErrNotFound:
			shift = store.NewShift(r.Name, n)
			shift.Start = start.Format(DateFormat)
			shift.End = end.Format(DateFormat)

		case err == nil:
			if shift.ShiftStatus != store.ShiftStatusScheduled {
				return nil, errors.Errorf("can not be scheduled, it is %q", shift.ShiftStatus)
			}
			if shift.Start != start.Format(DateFormat) || shift.End != end.Format(DateFormat) {
				return nil, errors.Errorf("loaded shift has wrong dates %v-%v, expected %v-%v",
					shift.Start, shift.End, start, end)
			}

		default:
			return nil, err
		}

		// Start with the users already scheduled in the shift, if any
		var shiftUsers store.UserList
		shiftUsers, err = api.loadUsers(shift.MattermostUserIDs)
		if err != nil {
			return nil, err
		}

		// replace with cached Users where appropriate
		for k := range shiftUsers {
			if cachedUsers[k] != nil {
				shiftUsers[k] = cachedUsers[k]
			}
		}

		err = api.prepareShift(r, rotationUsers, n, shift, shiftUsers, start, end, autofill)
		if err != nil {
			return nil, err
		}

		// Update shift's users' last served counter, and update the cache in
		// case they were not there
		for k, u := range shiftUsers {
			// rotationUSers is using the same pointers, is updated automatically
			u.Rotations[rotationName] = n
			cachedUsers[k] = u
		}

		shifts = append(shifts, shift)
	}

	return shifts, nil
}

func (api *api) prepareShift(r *store.Rotation, rotationUsers store.UserList, shiftNumber int, shift *store.Shift, shiftUsers store.UserList,
	start, end time.Time, autofill bool) error {

	// Filter out users who will not participate
	pool := store.UserList{}
	for _, u := range rotationUsers {
		// Do not include if already in the shift
		_, ok := shift.MattermostUserIDs[u.MattermostUserID]
		if ok {
			// logger.Debugf("skipping %v, already in the shift", u.MattermostUserID)
			continue
		}
		if !IsUserAvailable(u, start, end) {
			// logger.Debugf("skipping %v, unavailable", u.MattermostUserID)
			continue
		}
		pool[u.MattermostUserID] = u
	}

	var unsatisfied bool
	for _, need := range r.Needs {
		if len(shiftUsers) >= r.Size {
			unsatisfied = true
			break
		}
		need.Count = api.unsatisfiedNeed(need, shiftUsers)
		if need.Count == 0 {
			continue
		}

		if autofill {
			var picked store.UserList
			picked, err := api.pickUsersForNeed(r, need, pool, shiftNumber)
			if err != nil {
				return err
			}
			for _, u := range picked {
				shiftUsers[u.MattermostUserID] = u
			}
			need.Count = 0
		}
	}

	if unsatisfied {
		// TODO analyze failure
		return errors.New("<><> impossible")
	}

	// Backfill any remaining headcount from the remaining pool
	picked, err := api.pickUsersN(r, pool, shiftNumber, r.Size-len(shiftUsers))
	if err != nil {
		return err
	}
	for _, u := range picked {
		shiftUsers[u.MattermostUserID] = u
	}

	for _, u := range shiftUsers {
		shift.MattermostUserIDs[u.MattermostUserID] = u.MattermostUserID
	}

	return nil
}

func (api *api) unsatisfiedNeed(need store.Need, users store.UserList) int {
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

func (api *api) pickUsersForNeed(r *store.Rotation, need store.Need, users store.UserList, shiftNumber int) (store.UserList, error) {
	qualified := store.UserList{}
	for _, user := range users {
		skillLevel, _ := user.SkillLevels[need.Skill]
		if skillLevel >= need.Level {
			qualified[user.MattermostUserID] = user
		}
	}

	picked, err := api.pickUsersN(r, qualified, shiftNumber, need.Count)
	if err != nil {
		return nil, err
	}
	for k := range picked {
		delete(users, k)
	}
	return picked, nil
}

func (api *api) pickUsersN(r *store.Rotation, users store.UserList, shiftNumber int, numUsers int) (store.UserList, error) {
	picked := store.UserList{}
	for c := numUsers; c > 0; c-- {
		user, err := api.pickOne(r, users, shiftNumber)
		if err != nil {
			return nil, err
		}
		picked[user.MattermostUserID] = user
		delete(users, user.MattermostUserID)
	}

	return picked, nil
}

func (api *api) pickOne(r *store.Rotation, users store.UserList, shiftNumber int) (*store.User, error) {
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
	last, _ := user.Rotations[r.Name]
	return math.Pow(2.0, float64(shiftNumber-last))
}

func IsUserAvailable(user *store.User, start, end time.Time) bool {
	return true
}
