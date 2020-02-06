// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"sort"
	"time"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/autofill"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type fill struct {
	// Parameters
	bot.Logger
	rotationID  string
	size        int
	shiftNumber int
	userWeightF func(user *sl.User) float64

	// State
	pool             sl.UserMap
	chosen           sl.UserMap
	requiredNeeds    store.Needs
	needPools        map[string]sl.UserMap // uses skill-level for the key
	constrainedNeeds store.Needs
}

func makeAutofill(rotationID string, size int, needs store.Needs,
	pool sl.UserMap, chosen sl.UserMap, shiftNumber int, shiftStart, shiftEnd time.Time, logger bot.Logger) (*fill, error) {
	if chosen == nil {
		chosen = sl.UserMap{}
	}

	af := fill{
		Logger:        logger,
		rotationID:    rotationID,
		size:          size,
		pool:          pool,
		shiftNumber:   shiftNumber,
		requiredNeeds: store.Needs{},
		needPools:     map[string]sl.UserMap{},
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
			af.needPools[need.SkillLevel()] = sl.UsersQualifiedForNeed(af.pool, need)
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

func (af *fill) fill() (sl.UserMap, error) {
	af.Debugf(af.markdown())

	for len(af.chosen) < af.size {
		err := af.fillOne()
		if err != nil {
			return nil, err
		}
	}

	if len(af.requiredNeeds) > 0 {
		return nil, af.newError(nil, autofill.ErrSizeExceeded)
	}

	return af.chosen, nil
}

func (af *fill) fillOne() error {
	need, pool := af.pickNeed(af.requiredNeeds, af.needPools)
	if pool == nil {
		pool = af.pool
	}
	var user *sl.User
	for {
		user = af.pickUser(pool)
		if user == nil {
			if need != nil {
				return af.newError(need, autofill.ErrInsufficientForNeeds)
			} else {
				return af.newError(nil, autofill.ErrInsufficientForSize)
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

func (af *fill) meetsConstraints(user *sl.User) bool {
	for _, need := range af.constrainedNeeds {
		if sl.IsUserQualifiedForNeed(user, need) && need.Max-1 < 0 {
			af.Debugf("Disqualified %s against max on %s", user.Markdown(), need.Markdown())
			return false
		}
	}
	return true
}

func (af *fill) removeUser(user *sl.User) {
	delete(af.pool, user.MattermostUserID)
	for _, pool := range af.needPools {
		delete(pool, user.MattermostUserID)
	}
}

func (af *fill) acceptUser(user *sl.User) error {
	af.removeUser(user)

	// update the constraints
	for _, need := range af.constrainedNeeds {
		if sl.IsUserQualifiedForNeed(user, need) {
			need.Max--
		}
	}

	// update the requirements
	var updatedRequiredNeeds store.Needs
	for _, need := range af.requiredNeeds {
		pool := af.needPools[need.SkillLevel()]
		if sl.IsUserQualifiedForNeed(user, need) {
			need.Min--
			if need.Min == 0 {
				// filled to requirement, do not include in the updated map
				continue
			}
			if len(pool) == 0 {
				return af.newError(need, autofill.ErrInsufficientForNeeds)
			}
		}
		// the user is already removed from all needs' pools
		updatedRequiredNeeds = append(updatedRequiredNeeds, need)
	}
	af.requiredNeeds = updatedRequiredNeeds

	if af.chosen == nil {
		af.chosen = sl.UserMap{}
	}
	af.chosen[user.MattermostUserID] = user
	return nil
}

func (af *fill) markdown() string {
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

func (af *fill) newError(need *store.Need, err error) *autofill.Error {
	var unmet store.Needs
	for _, need := range af.requiredNeeds {
		unmet = append(unmet, need)
	}
	return &autofill.Error{
		Err:           err,
		UnmetNeeds:    unmet,
		UnmetNeed:     need,
		UnmetCapacity: af.size - len(af.chosen),
		ShiftNumber:   af.shiftNumber,
	}
}
