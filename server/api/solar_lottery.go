// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type autofillError struct {
	orig            error
	causeUnmetNeeds []store.Need
	causeUnmetNeed  *store.Need
	causeCapacity   int
	shiftNumber     int
}

func (e autofillError) Error() string {
	message := ""
	if e.causeCapacity > 0 {
		message = fmt.Sprintf("failed filling to capacity, missing %v", e.causeCapacity)
	}
	if e.causeUnmetNeed != nil {
		if message != "" {
			message += ", "
		}
		message += fmt.Sprintf("failed filling need %s", markdownNeed(*e.causeUnmetNeed))
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

func (api *api) autofillShift(rotation *Rotation, shiftNumber int, shift *Shift) error {
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
				api.Logger.Debugf("Disqualified user %s: unavailable", api.MarkdownUserWithSkills(user))
			}
		}
	}

	// Get a snapshot of the initial probability distribution
	var cdf *cdf
	var unmet []store.Need
	var pickedUser, chosenUser *User
	var qualifiedUsers UserMap
	cycle := func() {
		if pickedUser != nil {
			delete(pool, pickedUser.MattermostUserID)
			delete(qualifiedUsers, pickedUser.MattermostUserID)
		}
		if chosenUser != nil {
			chosen[chosenUser.MattermostUserID] = chosenUser
			chosenUser = nil
		}

		cdf = api.calculateCDF(rotation, pool, shiftNumber)
		rotation.sortNeedsByHeat(cdf)
		unmet = unmetNeeds(rotation.Needs, chosen)
	}

	// Set and snapshot the initial state for logging later
	cycle()

	out := fmt.Sprintf("filling %s, choosing from:\n", api.MarkdownShift(rotation, shiftNumber))
	sort.Sort(cdf)
	for i, id := range cdf.ids {
		out += fmt.Sprintf("- **%.3f**: %s\n", cdf.weights[i]/cdf.total, api.MarkdownUserWithSkills(pool[id]))
	}
	api.Logger.Debugf("%s", out)

CYCLE:
	for ; len(unmet) > 0; cycle() {
		need := unmet[0]
		autoerr := func(err error) autofillError {
			return autofillError{
				causeUnmetNeed:  &need,
				causeUnmetNeeds: unmet,
				causeCapacity:   rotation.Size - len(chosen),
				orig:            err,
				shiftNumber:     shiftNumber,
			}
		}

		if need.Max == 0 {
			// nothing else to do, not a realistic possibility unless misconfigured
			continue
		}

		if len(chosen) >= rotation.Size {
			return autoerr(ErrFailedSizeExceeded)
		}

		qualifiedUsers = usersQualifiedForNeed(pool, need)
		if len(qualifiedUsers) == 0 {
			return autoerr(ErrFailedInsufficientForNeeds)
		}

		for ; len(qualifiedUsers) > 0; cycle() {
			pickedUser = api.pickUser(rotation, qualifiedUsers, shiftNumber)
			if pickedUser == nil {
				// not really a possibility
				return autoerr(errors.Errorf("failed to pick a user out of %v qualified", len(qualifiedUsers)))
			}

			var maxedNeed *store.Need
			for _, rotationNeed := range rotation.Needs {
				if rotationNeed.Max < 0 {
					continue
				}
				if qualifiedForNeed(pickedUser, rotationNeed) &&
					rotationNeed.Max-len(usersQualifiedForNeed(chosen, rotationNeed))-1 < 0 {
					maxedNeed = &rotationNeed
					break
				}
			}

			if maxedNeed != nil {
				api.Logger.Debugf("Disqualified user %s from %s: would hit max on %s.",
					api.MarkdownUser(pickedUser), api.MarkdownNeed(need), api.MarkdownNeed(*maxedNeed))
				continue
			}

			chosenUser = pickedUser
			api.Logger.Debugf("Accepted user %s for %s", api.MarkdownUser(chosenUser), api.MarkdownNeed(need))
			continue CYCLE
		}
	}

	// All needs are met

	autoerr := func(err error) autofillError {
		return autofillError{
			causeCapacity: rotation.Size - len(chosen),
			orig:          err,
			shiftNumber:   shiftNumber,
		}
	}

	for ; len(chosen) < rotation.Size && len(pool) > 0; cycle() {
		// Choose next user from the pool, weighted-randomly.
		pickedUser = api.pickUser(rotation, pool, shiftNumber)
		if pickedUser == nil {
			// not really a possibility
			return autoerr(errors.Errorf("failed to pick a user out of %v qualified", len(qualifiedUsers)))
		}

		var maxedNeed *store.Need
		for _, rotationNeed := range rotation.Needs {
			if rotationNeed.Max < 0 {
				continue
			}
			if qualifiedForNeed(pickedUser, rotationNeed) &&
				rotationNeed.Max-len(usersQualifiedForNeed(chosen, rotationNeed))-1 < 0 {
				maxedNeed = &rotationNeed
				break
			}
		}

		if maxedNeed != nil {
			api.Logger.Debugf("Disqualified user %s from backfill: would hit max on %s.",
				api.MarkdownUser(pickedUser), api.MarkdownNeed(*maxedNeed))
		} else {
			chosenUser = pickedUser
			api.Logger.Debugf("Accepted user %s for backfill",
				api.MarkdownUser(chosenUser))
		}
	}

	if len(chosen) < rotation.Size {
		return autoerr(ErrFailedInsufficientForSize)
	}

	for _, user := range chosen {
		shift.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID
	}
	api.Logger.Debugf("Filled %s, chose %s. Details:\n%s",
		api.MarkdownShift(rotation, shiftNumber),
		api.MarkdownUsers(chosen),
		api.MarkdownShiftBullets(rotation, shiftNumber, shift))

	return nil
}

func qualifiedForNeed(user *User, need store.Need) bool {
	skillLevel, _ := user.SkillLevels[need.Skill]
	return skillLevel >= need.Level
}

func usersQualifiedForNeed(users UserMap, need store.Need) UserMap {
	qualified := UserMap{}
	for id, user := range users {
		if qualifiedForNeed(user, need) {
			qualified[id] = user
		}
	}
	return qualified
}

func (api *api) pickUser(rotation *Rotation, from UserMap, shiftNumber int) *User {
	cdf := api.calculateCDF(rotation, from, shiftNumber)
	random := rand.Float64() * cdf.total
	i := sort.Search(len(cdf.cdf), func(i int) bool {
		return cdf.cdf[i] >= random
	})
	if i < 0 || i >= len(cdf.cdf) {
		return nil
	}

	return from[cdf.ids[i]]
}

func unmetNeeds(needs []store.Need, users UserMap) []store.Need {
	work := append([]store.Need{}, needs...)
	for i, need := range work {
		for _, user := range users {
			if qualifiedForNeed(user, need) {
				work[i].Min--
				work[i].Max--
			}
		}
	}

	var unmet []store.Need
	for _, need := range work {
		if need.Min > 0 {
			unmet = append(unmet, need)
		}
	}
	return unmet
}

func unmetNeedsOld(needs []store.Need, users UserMap) []store.Need {
	unmet := []store.Need{}
	for _, need := range needs {
		cQualified := 0
		for _, user := range users {
			if qualifiedForNeed(user, need) {
				cQualified++
			}
		}
		if cQualified < need.Min {
			unmet = append(unmet, need)
		}
	}
	return unmet
}

func (rotation *Rotation) sortNeedsByHeat(cdf *cdf) {
	needWeights := make([]float64, len(rotation.Needs))
	for i, need := range rotation.Needs {
		if need.Min <= 0 {
			// Maximum weight will move the need to the bottom, so it's last matched against
			needWeights[i] = math.MaxFloat64
			continue
		}

		for j, id := range cdf.ids {
			user := cdf.users[id]
			if qualifiedForNeed(user, need) {
				needWeights[i] += cdf.weights[j]
			}
		}

		// Normalize per remaining needed user
		needWeights[i] = needWeights[i] / float64(need.Min)
	}

	sort.Sort(&weightedNeedSorter{
		needs:   rotation.Needs,
		weights: needWeights,
	})
}

type weightedNeedSorter struct {
	needs   []store.Need
	weights []float64
}

func (s *weightedNeedSorter) Len() int {
	return len(s.needs)
}

// Sort all needs that have a max limit to the top, to reduce hitting that
// limit opportunistically. Otherwise sort by normalized weight representing
// staffing heat level.
func (s *weightedNeedSorter) Less(i, j int) bool {
	if s.needs[i].Max >= 0 && s.needs[j].Max < 0 {
		return true
	}
	if s.needs[i].Max < 0 && s.needs[j].Max >= 0 {
		return false
	}
	return s.weights[i] < s.weights[j]
}

func (s *weightedNeedSorter) Swap(i, j int) {
	s.weights[i], s.weights[j] = s.weights[j], s.weights[i]
	s.needs[i], s.needs[j] = s.needs[j], s.needs[i]
}

type cdf struct {
	ids     []string
	weights []float64
	cdf     []float64
	users   UserMap
	total   float64
}

func (cdf *cdf) Len() int {
	return len(cdf.ids)
}

func (cdf *cdf) Less(i, j int) bool {
	return cdf.weights[i] > cdf.weights[j]
}

func (cdf *cdf) Swap(i, j int) {
	cdf.weights[i], cdf.weights[j] = cdf.weights[j], cdf.weights[i]
	cdf.ids[i], cdf.ids[j] = cdf.ids[j], cdf.ids[i]
	cdf.cdf[i], cdf.cdf[j] = cdf.cdf[j], cdf.cdf[i]
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
