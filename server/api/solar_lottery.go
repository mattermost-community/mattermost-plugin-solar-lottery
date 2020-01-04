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

type autofillError struct {
	orig                error
	unsatisfiedNeeds    []store.Need
	unsatisfiedCapacity int
	shiftNumber         int
}

func (e autofillError) Error() string {
	message := ""
	if e.unsatisfiedCapacity > 0 {
		message = fmt.Sprintf("couldn't fill shift %v to capacity, missing %v", e.shiftNumber, e.unsatisfiedCapacity)
	}
	if len(e.unsatisfiedNeeds) > 0 {
		if message != "" {
			message += ", "
		}
		message += fmt.Sprintf("couldn't fill shift %v needs %s", e.shiftNumber, MarkdownNeeds(e.unsatisfiedNeeds))
	}
	if e.orig != nil {
		return errors.WithMessage(e.orig, message).Error()
	} else {
		return message
	}
}

var ErrFailedInsufficientForNeeds = errors.New("failed to satisfy needs, not enough skilled users available")
var ErrFailedSizeExceeded = errors.New("failed to satisfy needs, exceeded rotation size")
var ErrFailedInsufficientForSize = errors.New("failed to satisfy rotation size requirement")

func (api *api) autofillShift(rotation *Rotation, shiftNumber int, shift *Shift, autofill bool) error {
	if len(shift.MattermostUserIDs) > 0 {
		api.Logger.Debugf("Shift %v already has users %v: %v",
			shiftNumber, len(shift.MattermostUserIDs), api.MarkdownUsersWithSkills(rotation.ShiftUsers(shift)))
	}

	pool := rotation.Users.Clone(false)
	chosen := rotation.ShiftUsers(shift)
	// remove any users already chosen from the pool
	for id := range chosen {
		_, ok := pool[id]
		if ok {
			delete(pool, id)
		}
	}
	// remove any unavailable users from the pool
	for _, user := range pool {
		overlappingEvents, err := user.overlapEvents(shift.StartTime, shift.EndTime, false)
		if err != nil {
			return autofillError{orig: err}
		}
		for _, event := range overlappingEvents {
			// Unavailable events apply to all rotations, Shift events apply
			//  only to the rotation from which they come.
			if event.Type == store.EventTypeUnavailable ||
				(event.Type == store.EventTypeShift && event.RotationID == rotation.RotationID) {

				delete(pool, user.MattermostUserID)
				api.Logger.Debugf("Unavailable user %s", api.MarkdownUserWithSkills(user))
			}
		}
	}

	// Calculate the initial probability distribution of the pool to log later.
	logPool := []string{}
	func() {
		cdf := api.calculateCDF(rotation, pool, shiftNumber)
		for i, id := range cdf.ids {
			logPool = append(logPool, fmt.Sprintf("%s: (%.2f)", api.MarkdownUser(pool[id]), cdf.weights[i]/cdf.total))
		}
	}()

	unsatisfiedNeeds := api.unsatisfiedNeeds(rotation.Needs, chosen)
	aErr := func(err error) error {
		return autofillError{
			unsatisfiedNeeds:    unsatisfiedNeeds,
			unsatisfiedCapacity: rotation.Size - len(chosen),
			orig:                err,
			shiftNumber:         shiftNumber,
		}
	}

	if !autofill {
		if len(unsatisfiedNeeds) > 0 {
			return aErr(errors.New("not autofilled"))
		}
		return nil
	}

	// First pass: start with the qualified, backfill the rest if any left
	unqualified := UserMap{}
	for len(pool) > 0 && len(unsatisfiedNeeds) > 0 {
		if rotation.Size != 0 && len(chosen) >= rotation.Size {
			// Reached capacity.
			break
		}

		user, cdf := api.pickUser(rotation, pool, shiftNumber)
		if user == nil {
			return aErr(ErrFailedInsufficientForNeeds)
		}

		if api.userIsQualifiedForShift(user, rotation, chosen, shift.StartTime, shift.EndTime, unsatisfiedNeeds) {
			chosen[user.MattermostUserID] = user
			delete(pool, user.MattermostUserID)
			unsatisfiedNeeds = api.unsatisfiedNeedsSorted(rotation, cdf, chosen)
			api.Logger.Debugf("Accepted user %s, shift: %v, pool: %v, %v unmet needs: %v", api.MarkdownUserWithSkills(user), len(chosen), len(pool), len(unsatisfiedNeeds), MarkdownNeeds(unsatisfiedNeeds))
		} else {
			unqualified[user.MattermostUserID] = user
			delete(pool, user.MattermostUserID)
			api.Logger.Debugf("Disqualified user %s, chosen: %v, pool: %v, %v unmet needs: %v", api.MarkdownUserWithSkills(user), len(chosen), len(pool), len(unsatisfiedNeeds), MarkdownNeeds(unsatisfiedNeeds))
		}
	}

	if len(unsatisfiedNeeds) > 0 {
		if len(pool) == 0 {
			return aErr(ErrFailedInsufficientForNeeds)
		} else {
			return aErr(ErrFailedSizeExceeded)
		}
	}

	// Second pass: backfill the rest from unqualified, merge any remaining in
	// the pool first.
	for id, user := range pool {
		unqualified[id] = user
	}

	for len(chosen) < rotation.Size && len(unqualified) > 0 {
		// Choose next user from the pool, weighted-randomly.
		user, _ := api.pickUser(rotation, unqualified, shiftNumber)
		if user == nil {
			return aErr(ErrFailedInsufficientForSize)
		}

		chosen[user.MattermostUserID] = user
		delete(unqualified, user.MattermostUserID)
		api.Logger.Debugf("Accepted user %s for backfill, chosen: %v, unqualified pool: %v", api.MarkdownUserWithSkills(user), len(chosen), len(unqualified))
	}

	if len(chosen) < rotation.Size {
		return aErr(ErrFailedInsufficientForSize)
	}

	for _, user := range chosen {
		shift.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID
	}
	api.Logger.Debugf("Shift %v (%v to %v) prepared, chose users %s **from** %s", shiftNumber, shift.Start, shift.End, api.MarkdownUsersWithSkills(chosen), logPool)

	return nil
}

// To simplify redundant checks, users that are currently serving in shift(s),
// qualify but will have 0 probability anyway so should never be chosen.
func (api *api) userIsQualifiedForShift(user *User, rotation *Rotation, ShiftUsers UserMap,
	start, end time.Time, unsatisfiedNeeds []store.Need) bool {

	// disqualify from any maxed out needs
	for _, need := range rotation.Needs {
		if need.Max <= 0 || !api.userIsQualifiedForNeed(user, need) {
			continue
		}

		if len(api.usersQualifiedForNeed(ShiftUsers, need)) >= need.Max {
			return false
		}
	}

	// see if qualifies
	if len(unsatisfiedNeeds) == 0 {
		return true
	}
	for _, need := range unsatisfiedNeeds {
		if !api.userIsQualifiedForNeed(user, need) {
			continue
		}
		return true
	}
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

func (api *api) pickUser(rotation *Rotation, pool UserMap, shiftNumber int) (*User, *cdf) {
	cdf := api.calculateCDF(rotation, pool, shiftNumber)
	random := rand.Float64() * cdf.total
	i := sort.Search(len(cdf.cdf), func(i int) bool {
		return cdf.cdf[i] >= random
	})
	if i < 0 || i >= len(cdf.cdf) {
		return nil, nil
	}

	return pool[cdf.ids[i]], cdf
}

func (api *api) unsatisfiedNeedsSorted(rotation *Rotation, poolCDF *cdf, chosen UserMap) []store.Need {
	unsatisfied := []store.Need{}
	unsatisfiedWeights := []float64{}
	for _, need := range rotation.Needs {
		cQualified := 0
		for _, user := range chosen {
			if api.userIsQualifiedForNeed(user, need) {
				cQualified++
			}
		}

		toFill := need.Min - cQualified
		if toFill > 0 {
			unsatisfied = append(unsatisfied, need)

			if poolCDF != nil {
				w := float64(0)
				for i, id := range poolCDF.ids {
					user := poolCDF.users[id]
					if api.userIsQualifiedForNeed(user, need) {
						w += poolCDF.weights[i]
					}
				}
				unsatisfiedWeights = append(unsatisfiedWeights, w/float64(toFill))
			}
		}
	}

	sort.Sort(&unsatisfiedNeedSorter{
		needs:   unsatisfied,
		weights: unsatisfiedWeights,
	})

	return unsatisfied
}

func (api *api) unsatisfiedNeeds(needs []store.Need, users UserMap) []store.Need {
	unsatisfied := []store.Need{}
	for _, need := range needs {
		cQualified := 0
		for _, user := range users {
			if api.userIsQualifiedForNeed(user, need) {
				cQualified++
			}
		}
		if cQualified < need.Min {
			unsatisfied = append(unsatisfied, need)
		}
	}

	return unsatisfied
}

type unsatisfiedNeedSorter struct {
	needs   []store.Need
	weights []float64
}

func (s *unsatisfiedNeedSorter) Len() int {
	return len(s.needs)
}
func (s *unsatisfiedNeedSorter) Less(i, j int) bool {
	return s.weights[i] < s.weights[j]
}
func (s *unsatisfiedNeedSorter) Swap(i, j int) {
	s.weights[i], s.weights[j] = s.weights[j], s.weights[i]
}

type cdf struct {
	ids     []string
	weights []float64
	cdf     []float64
	users   UserMap
	total   float64
}

func (api *api) calculateCDF(rotation *Rotation, users UserMap, shiftNumber int) *cdf {
	cdf := &cdf{
		ids:     []string{},
		weights: []float64{},
		users:   UserMap{},
	}

	for _, user := range users {
		weight := userWeight(rotation, user, shiftNumber)
		cdf.ids = append(cdf.ids, user.MattermostUserID)
		cdf.weights = append(cdf.weights, weight)
		cdf.total += weight
		cdf.users[user.MattermostUserID] = user
	}

	cdf.cdf = make([]float64, len(cdf.weights))
	floats.CumSum(cdf.cdf, cdf.weights)
	return cdf
}

func userWeight(rotation *Rotation, user *User, shiftNumber int) float64 {
	next, _ := user.LastServed[rotation.RotationID]
	if next > shiftNumber {
		// return a non-0 but very low number, to prevent user from serving
		// other than if all have 0 weights
		return 1e-12
	}
	return math.Pow(2.0, float64(shiftNumber-next))
}
