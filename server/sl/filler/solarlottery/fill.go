// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"sort"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type fill struct {
	bot.Logger
	userWeightF func(user *sl.User) float64

	// Parameters
	r              *sl.Rotation
	task           *sl.Task
	time           int64 // seconds, unix time
	doublingPeriod int64 // seconds

	// State
	pool         *sl.Users
	poolWeights  map[types.ID]float64
	filled       *sl.Users
	require      *sl.Needs
	requirePools map[types.ID]*sl.Users // by need ID (SkillLevel as string)
	limit        *sl.Needs
}

func newFill(r *sl.Rotation, t *sl.Task, forTime types.Time, logger bot.Logger) *fill {
	if forTime.IsZero() && !t.ExpectedStart.IsZero() {
		forTime = t.ExpectedStart
	}
	if forTime.IsZero() {
		forTime = types.NewTime(time.Now())
	}

	pool := sl.NewUsers()
	if r.Users != nil {
		pool = r.Users.Clone()
	}

	f := fill{
		Logger:         logger,
		r:              r,
		task:           t,
		time:           forTime.Unix(),
		pool:           pool,
		poolWeights:    map[types.ID]float64{},
		filled:         sl.NewUsers(),
		require:        t.Require.Clone(),
		limit:          t.Limit.Clone(),
		requirePools:   map[types.ID]*sl.Users{},
		doublingPeriod: 7 * 24 * 3600, // 14 days
	}
	f.userWeightF = f.userWeight

	// remove any unavailable users from the pool
	for _, user := range f.pool.AsArray() {
		overlapping := user.FindUnavailable(
			types.NewDurationInterval(t.ExpectedStart, t.ExpectedDuration), r.RotationID)
		if len(overlapping) == 0 {
			continue
		}
		f.pool.Delete(user.MattermostUserID)
		logger.Debugf("Disqualified %s: unavailable", user.Markdown())
	}

	// fill in all users already in the task
	for _, user := range t.Users.AsArray() {
		_ = f.fillUser(user, true)
		f.Debugf("%s is already assigned", user.MarkdownWithSkills())
	}
	f.trimRequire()

	// create pools for all required needs
	for _, need := range f.require.AsArray() {
		qualified, _ := need.QualifyUsers(f.pool)
		f.requirePools[need.GetID()] = qualified
	}

	return &f
}

func (f *fill) fill() (*sl.Users, error) {
	f.Debugf(f.markdown())

	for {
		done, need := f.pickRequiredNeed()
		if done {
			break
		}

		for {
			user := f.pickUser(f.requirePools[need.GetID()])
			if user == nil {
				return nil, f.newError(need, sl.ErrFillInsufficient)
			}

			// The picked user is either accepted, or declined based on Limit
			// constraints, so remove it from the pools right away
			violated := f.fillUser(user, false)
			if !violated.IsEmpty() {
				f.Debugf("...skipped user %s: would exceed limits on %s", user.Markdown(), violated.Markdown())
				continue
			}
			f.Debugf("...picked %s for %s, remaining R: %s, L: %s", user.MarkdownWithSkills(), need.Markdown(), f.require.Markdown(), f.limit.Markdown())
			break
		}
	}

	f.Debugf("filled %s into %s", f.filled.MarkdownWithSkills(), f.task.Markdown())

	return f.filled, nil
}

func (f *fill) fillUser(user *sl.User, preassigned bool) (violated *sl.Needs) {
	// The picked user is either accepted, or declined based on Limit
	// constraints, so remove it from the pools right away
	f.pool.Delete(user.MattermostUserID)
	for _, pool := range f.requirePools {
		pool.Delete(user.MattermostUserID)
	}

	updatedLimit, _, violated := f.limit.CheckLimits(user)
	if !preassigned && !violated.IsEmpty() {
		return violated
	}
	updatedRequire := f.require.CheckRequired(user)

	f.limit = updatedLimit
	f.require = updatedRequire
	if !preassigned {
		f.filled.Set(user)
	}
	return violated
}

func (f *fill) markdown() string {
	ws := newSorter(0)
	total := float64(0)
	for _, user := range f.pool.AsArray() {
		weight := f.userWeight(user)
		ws.Append(user.MattermostUserID, weight)
		total += weight
	}
	sort.Sort(ws)
	out := ""
	out += fmt.Sprintf("filling task %s:\n", f.task.Markdown())
	if !f.require.IsEmpty() {
		out += fmt.Sprintf("- Requires: %s\n", f.require.Markdown())
	}
	if !f.limit.IsEmpty() {
		out += fmt.Sprintf("- Limits: %s\n", f.limit.Markdown())
	}
	out += fmt.Sprintf("- User pool:\n")
	for i, id := range ws.IDs {
		user := f.pool.Get(id)
		out += fmt.Sprintf("  1. **%.3f**: %s\n", ws.Weights[i]/total, user.MarkdownWithSkills())
	}
	return out
}

func (f *fill) newError(need sl.Need, err error) *sl.FillError {
	unmet := sl.NewNeeds()
	for _, need := range f.require.AsArray() {
		unmet.Set(need)
	}
	return &sl.FillError{
		Err:        err,
		UnmetNeeds: unmet,
		FailedNeed: &need,
		TaskID:     f.task.TaskID,
	}
}
