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

var ErrInsufficientForNeeds = errors.New("failed to satisfy needs, not enough skilled users available")
var ErrSizeExceeded = errors.New("failed to satisfy needs, exceeded rotation size")
var ErrInsufficientForSize = errors.New("failed to satisfy rotation size requirement")

type autofill struct {
	api API

	size        int
	shiftNumber int

	pool   UserMap
	chosen UserMap

	// requiredNeeds contains the map of all needs that are still required, to
	// the pool of available, qualified, not yet chosen users for each need.
	requiredNeeds map[*store.Need]UserMap

	constrainedNeeds []*store.Need
}

func (api *api) autofillShift(rotation *Rotation, shiftNumber int, shift *Shift) error {
	if len(shift.MattermostUserIDs) > 0 {
		api.Logger.Debugf("Shift %v already has users %v: %v",
			shiftNumber, len(shift.MattermostUserIDs), api.MarkdownUsersWithSkills(rotation.ShiftUsers(shift)))
	}

	af, err := api.makeAutofill(rotation, shiftNumber, shift)
	if err != nil {
		return err
	}
	api.Logger.Debugf("Autofill: %s", af.markdown(rotation, shiftNumber))

	chosen, err := af.fill()
	if err != nil {
		return err
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

func (af *autofill) fill() (UserMap, error) {
	for len(af.chosen) < af.size {
		err := af.fillOne()
		if err != nil {
			return nil, err
		}
	}

	if len(af.requiredNeeds) > 0 {
		return nil, af.newError(nil, ErrSizeExceeded)
	}

	return af.chosen, nil
}

func (af *autofill) fillOne() error {
	need, pool := af.hottestRequiredNeed()
	if pool == nil {
		pool = af.pool
	}

	var user *User
	for {
		user = af.pickUser(pool)
		if user == nil {
			if need != nil {
				return af.newError(need, ErrInsufficientForNeeds)
			} else {
				return af.newError(nil, ErrInsufficientForSize)
			}
		}

		if af.meetsConstraints(user) {
			break
		}

		af.dismissUser(user)
	}

	return af.acceptUser(user)
}

func (af *autofill) meetsConstraints(user *User) bool {
	for _, need := range af.constrainedNeeds {
		if qualifiedForNeed(user, *need) && need.Max-1 < 0 {
			return false
		}
	}
	return true
}

func (af *autofill) dismissUser(user *User) {
	delete(af.pool, user.MattermostUserID)
	for _, pool := range af.requiredNeeds {
		delete(pool, user.MattermostUserID)
	}
}

func (af *autofill) acceptUser(user *User) error {
	// remove the user from all pools whether it is accepted or not.
	af.dismissUser(user)

	// update the constraints
	for _, need := range af.constrainedNeeds {
		if qualifiedForNeed(user, *need) {
			need.Max--
		}
	}

	// update the requirements
	updatedRequiredNeeds := map[*store.Need]UserMap{}
	for need, pool := range af.requiredNeeds {
		if qualifiedForNeed(user, *need) {
			need.Min--
			if need.Min == 0 {
				// filled to requirement, do not include in the updated map
				continue
			}
			if len(pool) == 0 {
				return af.newError(need, ErrInsufficientForNeeds)
			}
			// the user is already removed from the pool
			updatedRequiredNeeds[need] = pool
		}
	}
	af.requiredNeeds = updatedRequiredNeeds

	af.chosen[user.MattermostUserID] = user
	return nil
}

func (af *autofill) pickUser(from UserMap) *User {
	if len(from) == 0 {
		return nil
	}

	cdf := make([]float64, len(from))
	weights := []float64{}
	ids := []string{}
	total := float64(0)
	for id, user := range from {
		ids = append(ids, id)
		weights = append(weights, user.weight)
		total += user.weight
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

func userWeight(rotationID string, user *User, shiftNumber int) float64 {
	last, _ := user.LastServed[rotationID]
	if last > shiftNumber {
		// return a non-0 but very low number, to prevent user from serving
		// other than if all have 0 weights
		return 1e-12
	}
	return math.Pow(2.0, float64(shiftNumber-last))
}

func (af *autofill) hottestRequiredNeed() (*store.Need, UserMap) {
	if len(af.requiredNeeds) == 0 {
		return nil, nil
	}

	s := &weightedNeedSorter{
		weights: make([]float64, len(af.requiredNeeds)),
		needs:   make([]*store.Need, len(af.requiredNeeds)),
	}
	i := 0
	for need, pool := range af.requiredNeeds {
		for _, user := range pool {
			s.weights[i] += user.weight
		}

		// Normalize per remaining needed user
		s.weights[i] = s.weights[i] / float64(need.Min)
		i++
	}

	sort.Sort(s)
	need := s.needs[0]
	pool := af.requiredNeeds[need]
	return need, pool
}

func (api *api) makeAutofill(rotation *Rotation, shiftNumber int, shift *Shift) (autofill, error) {
	af := autofill{
		api:           api,
		size:          rotation.Size,
		pool:          rotation.Users.Clone(false),
		chosen:        rotation.ShiftUsers(shift),
		shiftNumber:   shiftNumber,
		requiredNeeds: map[*store.Need]UserMap{},
	}

	// remove any users already chosen from the pool
	for id := range af.chosen {
		_, ok := af.pool[id]
		if ok {
			delete(af.pool, id)
		}
	}

	// remove any unavailable users from the pool, update weights
	for _, user := range af.pool {
		overlappingEvents, err := user.overlapEvents(shift.StartTime, shift.EndTime, false)
		if err != nil {
			return autofill{}, autofillError{orig: err}
		}
		for _, event := range overlappingEvents {
			// Unavailable events apply to all rotations, Shift events apply
			//  only to the rotation from which they come.
			if event.Type == store.EventTypePersonal ||
				(event.Type == store.EventTypeShift && event.RotationID == rotation.RotationID) {

				delete(af.pool, user.MattermostUserID)
				api.Logger.Debugf("Disqualified user %s: unavailable", api.MarkdownUserWithSkills(user))
			}
		}

		user.weight = userWeight(rotation.RotationID, user, shiftNumber)
	}

	// sort out the need requirements and constraints
	for _, need := range rotation.Needs {
		if need.Min > 0 {
			af.requiredNeeds[&need] = usersQualifiedForNeed(af.pool, need)
		}
		if need.Max >= 0 {
			af.constrainedNeeds = append(af.constrainedNeeds, &need)
		}
	}

	return af, nil
}

func (af *autofill) markdown(rotation *Rotation, shiftNumber int) string {
	ws := weightedUserSorter{}
	total := float64(0)
	for id, user := range af.pool {
		ws.ids = append(ws.ids, id)
		ws.weights = append(ws.weights, user.weight)
		total += user.weight
	}
	sort.Sort(&ws)
	out := fmt.Sprintf("filling %s, choosing from:\n", af.api.MarkdownShift(rotation, shiftNumber))
	for i, id := range ws.ids {
		out += fmt.Sprintf("- **%.3f**: %s\n", ws.weights[i]/total, af.api.MarkdownUserWithSkills(af.pool[id]))
	}
	return out
}

type autofillError struct {
	orig            error
	causeUnmetNeeds []store.Need
	causeUnmetNeed  *store.Need
	causeCapacity   int
	shiftNumber     int
}

func (af *autofill) newError(need *store.Need, err error) autofillError {
	unmet := []store.Need{}
	for need := range af.requiredNeeds {
		unmet = append(unmet, *need)
	}
	return autofillError{
		orig:            err,
		causeUnmetNeeds: unmet,
		causeUnmetNeed:  need,
		causeCapacity:   af.size - len(af.chosen),
		shiftNumber:     af.shiftNumber,
	}
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

type weightedNeedSorter struct {
	needs   []*store.Need
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
