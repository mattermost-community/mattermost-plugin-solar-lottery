// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrShiftAlreadyExists = errors.New("Shift already exists")

func (api *api) prepareShift(r *store.Rotation, shiftNumber int) (*store.Shift, error) {
	logger := api.Logger.With(bot.LogContext{
		"RotationName": r.Name,
		"ShiftNumber":  shiftNumber,
	})
	start, end, err := ShiftDates(r, shiftNumber)
	if err != nil {
		return nil, err
	}

	shift, err := api.ShiftStore.LoadShift(r.Name, shiftNumber)
	switch err {
	case store.ErrNotFound:
		shift = store.NewShift(r.Name, shiftNumber)
		shift.Start = start.Format(DateFormat)
		shift.End = end.Format(DateFormat)
		logger.Debugf("No existing shift, created: %v to %v", shift.Start, shift.End)

	case nil:
		if shift.ShiftStatus != store.ShiftStatusScheduled {
			return nil, errors.Errorf("can not be scheduled, it is %q", shift.ShiftStatus)
		}
		if shift.Start != start.Format(DateFormat) || shift.End != end.Format(DateFormat) {
			return nil, errors.Errorf("loaded shift has wrong dates %v-%v, expected %v-%v",
				shift.Start, shift.End, start, end)
		}
		logger.Debugf("Loaded existing shift: %v to %v", shift.Start, shift.End)

	default:
		return nil, err
	}

	// Load all rotation users
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
			logger.Debugf("skipping %v, already in the shift", u.MattermostUserID)
			continue
		}
		if !IsUserAvailable(u, start, end) {
			logger.Debugf("skipping %v, already in the shift", u.MattermostUserID)
			continue
		}
		pool[u.MattermostUserID] = u
	}
	logger.Debugf("%v users in the pool", len(pool))

	var unsatisfied bool
	// Start with the users already scheduled in the shift, if any
	selectedUsers, err := api.loadUsers(shift.MattermostUserIDs)
	if err != nil {
		return nil, err
	}
	logger.Debugf("%v users already in the shift", len(selectedUsers))

	for _, need := range r.Needs {
		needName := fmt.Sprintf("%s (%s)", need.Skill, Level(need.Level))
		logger.Debugf("preparing need %s", needName)
		if len(selectedUsers) >= r.Size {
			logger.Debugf("%v shift size reached, %v", r.Size, len(selectedUsers))
			unsatisfied = true
			break
		}
		need.Count = api.unsatisfiedNeed(need, selectedUsers)
		if need.Count == 0 {
			logger.Debugf("%s need already satisfied", needName)
			continue
		}

		picked, err := api.pickUsers(r, need, pool, shiftNumber)
		if err != nil {
			return nil, err
		}
		logger.Debugf("picked %v users, needed %v for %s", len(picked), need.Count, needName)

		for _, u := range picked {
			selectedUsers[u.MattermostUserID] = u
		}
		need.Count = 0
	}

	if unsatisfied {
		// TODO analyze failure
		return nil, errors.New("<><> impossible")
	}

	for _, u := range selectedUsers {
		shift.MattermostUserIDs[u.MattermostUserID] = u.MattermostUserID
	}

	return shift, nil
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

func (api *api) pickUsers(r *store.Rotation, need store.Need, users store.UserList, shiftNumber int) (store.UserList, error) {
	skilled := store.UserList{}
	for _, u := range users {
		skillLevel, _ := u.SkillLevels[need.Skill]
		if skillLevel >= need.Level {
			skilled[u.MattermostUserID] = u
			api.Logger.Debugf("Added %q (%s) to skilled for %s(%s)",
				u.MattermostUserID, Level(u.SkillLevels[need.Skill]), need.Skill, Level(need.Level))
		}
	}

	picked := store.UserList{}
	for c := need.Count; c > 0; c-- {
		user, err := api.pickOne(r, need, skilled, shiftNumber)
		if err != nil {
			return nil, err
		}
		picked[user.MattermostUserID] = user
	}

	return picked, nil
}

func (api *api) pickOne(r *store.Rotation, need store.Need, users store.UserList, shiftNumber int) (*store.User, error) {
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
