// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"gonum.org/v1/gonum/floats"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type autofill struct {
	// Parameters
	bot.Logger
	rotationID  string
	size        int
	shiftNumber int
	userWeightF func(user *api.User) float64

	// State
	pool             api.UserMap
	chosen           api.UserMap
	requiredNeeds    store.Needs
	needPools        map[string]api.UserMap // uses skill-level for the key
	constrainedNeeds store.Needs
}

func makeAutofill(rotationID string, size int, needs store.Needs,
	pool api.UserMap, chosen api.UserMap, shiftNumber int, shiftStart, shiftEnd time.Time, logger bot.Logger) (*autofill, error) {
	if chosen == nil {
		chosen = api.UserMap{}
	}

	af := autofill{
		Logger:        logger,
		rotationID:    rotationID,
		size:          size,
		pool:          pool,
		shiftNumber:   shiftNumber,
		requiredNeeds: store.Needs{},
		needPools:     map[string]api.UserMap{},
	}
	af.userWeightF = af.userWeight

	// remove any unavailable users from the pool, update weights
	for _, user := range af.pool {
		overlappingEvents, err := user.OverlapEvents(shiftStart, shiftEnd, false)
		if err != nil {
			return nil, err
		}
		for _, event := range overlappingEvents {
			// Unavailable events apply to all rotations, Shift events apply
			// only to the rotation from which they come.
			if event.Type == store.EventTypePersonal ||
				(event.Type == store.EventTypeShift && event.RotationID == rotationID) {

				delete(af.pool, user.MattermostUserID)
				logger.Debugf("Disqualified %s: unavailable", user.Markdown())
			}
		}
	}

	// sort out the need requirements and constraints
	for _, need := range needs {
		if need.Min > 0 {
			af.requiredNeeds = append(af.requiredNeeds, need)
			af.needPools[need.SkillLevel()] = api.UsersQualifiedForNeed(af.pool, need)
		}
		if need.Max >= 0 {
			af.constrainedNeeds = append(af.constrainedNeeds, need)
		}
	}

	for _, u := range chosen {
		af.acceptUser(u)
	}

	return &af, nil
}

func (af *autofill) fill() (api.UserMap, error) {
	af.Debugf(af.markdown())
	for len(af.chosen) < af.size {
		err := af.fillOne()
		if err != nil {
			return nil, err
		}
	}

	if len(af.requiredNeeds) > 0 {
		return nil, af.newError(nil, api.ErrSizeExceeded)
	}

	return af.chosen, nil
}

func (af *autofill) fillOne() error {
	need, pool := af.hottestRequiredNeed(af.requiredNeeds, af.needPools)
	if pool == nil {
		pool = af.pool
	}
	var user *api.User
	for {
		user = af.pickUser(pool)
		if user == nil {
			if need != nil {
				return af.newError(need, api.ErrInsufficientForNeeds)
			} else {
				return af.newError(nil, api.ErrInsufficientForSize)
			}
		}

		if af.meetsConstraints(user) {
			break
		}

		// Dismiss
		af.removeUser(user)
	}

	return af.acceptUser(user)
}

func (af *autofill) meetsConstraints(user *api.User) bool {
	for _, need := range af.constrainedNeeds {
		if api.IsUserQualifiedForNeed(user, need) && need.Max-1 < 0 {
			af.Debugf("Disqualified %s against max on %s", user.Markdown(), need.Markdown())
			return false
		}
	}
	return true
}

func (af *autofill) removeUser(user *api.User) {
	delete(af.pool, user.MattermostUserID)
	for _, pool := range af.needPools {
		delete(pool, user.MattermostUserID)
	}
}

func (af *autofill) acceptUser(user *api.User) error {
	af.removeUser(user)

	// update the constraints
	for _, need := range af.constrainedNeeds {
		if api.IsUserQualifiedForNeed(user, need) {
			need.Max--
		}
	}

	// update the requirements
	var updatedRequiredNeeds store.Needs
	for _, need := range af.requiredNeeds {
		pool := af.needPools[need.SkillLevel()]
		if api.IsUserQualifiedForNeed(user, need) {
			need.Min--
			if need.Min == 0 {
				// filled to requirement, do not include in the updated map
				continue
			}
			if len(pool) == 0 {
				return af.newError(need, api.ErrInsufficientForNeeds)
			}
		}
		// the user is already removed from all needs' pools
		updatedRequiredNeeds = append(updatedRequiredNeeds, need)
	}
	af.requiredNeeds = updatedRequiredNeeds

	if af.chosen == nil {
		af.chosen = api.UserMap{}
	}
	af.chosen[user.MattermostUserID] = user
	return nil
}

func (af *autofill) pickUser(from api.UserMap) *api.User {
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

func (af *autofill) userWeight(user *api.User) float64 {
	lastServed := user.LastServed[af.rotationID]
	if lastServed > af.shiftNumber {
		// return a non-0 but very low number, to prevent user from serving
		// other than if all have 0 weights
		return 1e-12
	}
	return math.Pow(2.0, float64(af.shiftNumber-lastServed))
}

func (af *autofill) hottestRequiredNeed(requiredNeeds store.Needs, needPools map[string]api.UserMap) (*store.Need, api.UserMap) {
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

func (af *autofill) markdown() string {
	ws := weightedUserSorter{}
	total := float64(0)
	for id, user := range af.pool {
		ws.ids = append(ws.ids, id)
		weight := af.userWeight(user)
		ws.weights = append(ws.weights, weight)
		total += weight
	}
	sort.Sort(&ws)
	out := ""
	out += fmt.Sprintf("filling shift %v, choosing from:\n", af.shiftNumber)
	for i, id := range ws.ids {
		out += fmt.Sprintf("- **%.3f**: %s\n", ws.weights[i]/total, af.pool[id].MarkdownWithSkills())
	}
	return out
}

func (af *autofill) newError(need *store.Need, err error) *api.AutofillError {
	var unmet store.Needs
	for _, need := range af.requiredNeeds {
		unmet = append(unmet, need)
	}
	return &api.AutofillError{
		Err:           err,
		UnmetNeeds:    unmet,
		UnmetNeed:     need,
		UnmetCapacity: af.size - len(af.chosen),
		ShiftNumber:   af.shiftNumber,
	}
}
