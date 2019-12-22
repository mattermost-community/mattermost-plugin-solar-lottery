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
)

func (api *api) loadOrMakeOneShift(rotation *Rotation, shiftNumber int, autofill bool) (*Shift, bool, error) {
	start, end, err := rotation.ShiftDatesForNumber(shiftNumber)
	if err != nil {
		return nil, false, err
	}

	var shift *Shift
	created := false
	storedShift, err := api.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	switch err {
	case nil:
		shift = &Shift{
			Shift: storedShift,
		}
		err = api.expandShift(shift)

	case store.ErrNotFound:
		if !autofill {
			return nil, false, err
		}
		shift, err = rotation.makeShift(shiftNumber, nil)
		if err != nil {
			return nil, false, err
		}
		created = true

	default:
		return nil, false, err
	}

	if shift.Start != start.Format(DateFormat) || shift.End != end.Format(DateFormat) {
		return nil, false, errors.Errorf("loaded shift has wrong dates %v-%v, expected %v-%v",
			shift.Start, shift.End, start, end)
	}

	err = api.expandShift(shift)
	if err != nil {
		return nil, false, err
	}

	return shift, created, nil
}

func (api *api) autofillShift(rotation *Rotation, shiftNumber int, shift *Shift, autofill bool) error {
	if len(shift.Users) > 0 {
		api.Logger.Debugf("Shift %v already has users %v: %v",
			shiftNumber, len(shift.Users), MarkdownUserMapWithSkills(shift.Users))
	}

	pool := rotation.Users.Clone()
	chosen := shift.Users.Clone()
	// remove any users already chosen from the pool
	for id := range chosen {
		_, ok := pool[id]
		if ok {
			delete(pool, id)
		}
	}

	unsatisfiedNeeds := api.unsatisfiedNeeds(rotation.Needs, chosen)

	if !autofill {
		if len(unsatisfiedNeeds) > 0 {
			return errors.Errorf("%v needs could not be satisfied, %v users in the shift, not autofilled.",
				len(unsatisfiedNeeds), len(chosen))
		}
		return nil
	}

	// Take a snapshot of the initial probability distribution for logging later
	logIDs, logCDF := cdf(rotation, pool, shiftNumber)
	logPool := []string{}
	logChosen := []string{}
	prev := float64(0)
	for i, id := range logIDs {
		logPool = append(logPool,
			fmt.Sprintf("%s: (%.2f)", MarkdownUser(pool[id]), (logCDF[i]-prev)/logCDF[len(logCDF)-1]))
		prev = logCDF[i]
	}

	// First pass: start with the qualified, backfill the rest if any left
	unqualified := UserMap{}
	for len(pool) > 0 {
		if rotation.Size != 0 && len(chosen) >= rotation.Size {
			api.Logger.Debugf("Reached capacity `%v`.", rotation.Size)
			break
		}
		if rotation.Size == 0 && len(unsatisfiedNeeds) == 0 {
			api.Logger.Debugf("Met all needs.")
			break
		}

		// Choose next user from the pool, weighted-randomly.
		ids, cdf := cdf(rotation, pool, shiftNumber)
		random := rand.Float64() * cdf[len(cdf)-1]
		i := sort.Search(len(cdf), func(i int) bool {
			return cdf[i] > random
		})
		user := pool[ids[i]]

		if api.userIsQualifiedForShift(user, rotation, chosen, shift.StartTime, shift.EndTime, unsatisfiedNeeds) {
			api.Logger.Debugf("`%v` needs remain, `%v` users in the pool, user %v ACCEPTED",
				len(unsatisfiedNeeds), len(pool), usernameWithSkills(user))

			chosen[user.MattermostUserID] = user
			logChosen = append(logChosen, usernameWithSkills(user))
			delete(pool, user.MattermostUserID)
		} else {
			// api.Logger.Debugf("%v needs remain, %v users in the pool, user %v rejected.",
			// 	len(unsatisfiedNeeds), len(pool), usernameWithSkills(user))
			unqualified[user.MattermostUserID] = user
		}

		unsatisfiedNeeds = api.unsatisfiedNeeds(rotation.Needs, chosen)
	}

	if len(unsatisfiedNeeds) > 0 {
		if len(pool) == 0 {
			return errors.Errorf("%v needs could not be satisfied", len(unsatisfiedNeeds))
		} else {
			return errors.Errorf("reached rotation capacity %v before %v needs could be satisfied",
				rotation.Size, len(unsatisfiedNeeds))
		}
	}

	// Second pass: backfill the rest from any unqualified
	for len(chosen) < rotation.Size && len(unqualified) > 0 {
		// Choose next user from the pool, weighted-randomly.
		ids, cdf := cdf(rotation, unqualified, shiftNumber)
		random := rand.Float64() * cdf[len(cdf)-1]
		i := sort.Search(len(cdf), func(i int) bool {
			return cdf[i] > random
		})
		user := unqualified[ids[i]]

		api.Logger.Debugf("backfilling `%v` unqualified, `%v` users in the pool, user %s ACCEPTED",
			rotation.Size-len(chosen), len(chosen), usernameWithSkills(user))

		chosen[user.MattermostUserID] = user
		logChosen = append(logChosen, usernameWithSkills(user))
		delete(unqualified, user.MattermostUserID)
	}

	if len(chosen) < rotation.Size {
		return errors.Errorf("%v unqualified users coulf not be filled", rotation.Size-len(chosen))
	}

	for _, user := range chosen {
		shift.Users[user.MattermostUserID] = user
		shift.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID
	}
	api.Logger.Debugf("Shift %v (%v to %v) prepared, chose **%v** users %s **from** %s", shiftNumber, shift.Start, shift.End, len(logChosen), logChosen, logPool)

	return nil
}

func cdf(rotation *Rotation, users UserMap, shiftNumber int) (ids []string, cdf []float64) {
	var weights []float64
	for _, user := range users {
		ids = append(ids, user.MattermostUserID)
		weights = append(weights, userWeight(rotation, user, shiftNumber))
	}

	cdf = make([]float64, len(weights))
	floats.CumSum(cdf, weights)
	return ids, cdf
}

func userWeight(rotation *Rotation, user *User, shiftNumber int) float64 {
	last, _ := user.Rotations[rotation.RotationID]
	return math.Pow(2.0, float64(shiftNumber-last))
}

func userIsAvailable(user *User, start, end time.Time) bool {
	//TODO userIsAvailable
	return true
}

// To simplify redundant checks, users that are currently serving in shift(s),
// qualify but will have 0 probability anyway so should never be chosen.
func (api *api) userIsQualifiedForShift(user *User, rotation *Rotation, shiftUsers UserMap,
	start, end time.Time, unsatisfiedNeeds map[string]store.Need) bool {

	if !userIsAvailable(user, start, end) {
		api.Logger.Debugf("DISQUALIFY user %v: NOT AVAILABLE", usernameWithSkills(user))
		return false
	}

	// disqualify from any maxed out needs
	for needName, need := range rotation.Needs {
		if need.Max <= 0 || !api.userIsQualifiedForNeed(user, need) {
			continue
		}

		if len(api.usersQualifiedForNeed(shiftUsers, need)) >= need.Max {
			api.Logger.Debugf("DISQUALIFY user %s, would exceed the max `%v` on `%s`",
				usernameWithSkills(user), need.Max, needName)
			return false
		}
	}

	// see if qualifies
	if len(unsatisfiedNeeds) == 0 {
		api.Logger.Debugf("QUALIFY user %v for NO NEED", usernameWithSkills(user))
		return true
	}
	for _, need := range unsatisfiedNeeds {
		if !api.userIsQualifiedForNeed(user, need) {
			continue
		}
		return true
	}
	api.Logger.Debugf("DISQUALIFY user %s skills did not qualify.", usernameWithSkills(user))
	return false
}

func (api *api) userIsQualifiedForNeed(user *User, need store.Need) bool {
	skillLevel, _ := user.SkillLevels[need.Skill]
	return skillLevel >= need.Level
}

func (api *api) usersQualifiedForNeed(users UserMap, need store.Need) UserMap {
	qualified := UserMap{}
	for id, user := range users {
		if api.userIsQualifiedForNeed(user, need) {
			qualified[id] = user
		}
	}
	return qualified
}

func (api *api) unsatisfiedNeeds(needs map[string]store.Need, users UserMap) map[string]store.Need {
	unsatisfied := map[string]store.Need{}
	for name, need := range needs {
		cQualified := 0
		for _, user := range users {
			if api.userIsQualifiedForNeed(user, need) {
				cQualified++
			}
		}
		if cQualified < need.Min {
			unsatisfied[name] = need
		}
	}
	return unsatisfied
}

func usernameWithSkills(user *User) string {
	return fmt.Sprintf("%s %s", MarkdownUser(user), MarkdownUserSkills(user))
}
