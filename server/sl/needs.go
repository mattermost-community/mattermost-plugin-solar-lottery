// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

// Needs is a map of SkillLevel to an int64 headcount needed.
type Needs struct {
	types.IntSet
}

var NeedOneAnyLevel = NewNeed(1, AnySkillLevel)

func NewNeeds(nn ...Need) *Needs {
	needs := &Needs{
		IntSet: *types.NewIntSet(),
	}
	for _, need := range nn {
		needs.Set(need)
	}
	return needs
}

func (needs Needs) Clone() *Needs {
	return NewNeeds(needs.AsArray()...)
}

func (needs *Needs) IsEmpty() bool {
	return needs == nil || needs.IntSet.IsEmpty()
}

func (needs Needs) Get(id types.ID) Need {
	if !needs.Contains(id) {
		return NewNeed(0, AnySkillLevel)
	}
	return NewNeed(needs.IntSet.Get(id), ParseSkillLevel(id))
}

func (needs Needs) GetCountForSkillLevel(skillLevel SkillLevel) int64 {
	id := types.ID(skillLevel.String())
	return needs.IntSet.Get(id)
}

func (needs *Needs) Set(need Need) {
	needs.IntSet.Set(need.ID, need.Value)
}

func (needs *Needs) SetCountForSkillLevel(skillLevel SkillLevel, count int64) {
	id := types.ID(skillLevel.String())
	needs.IntSet.Set(id, count)
}

func (needs Needs) Markdown() string {
	ss := []string{}
	for _, need := range needs.AsArray() {
		ss = append(ss, fmt.Sprintf("**%v** %s", need.Count(), need.SkillLevel()))
	}
	return strings.Join(ss, ", ")
}

func (needs Needs) MarkdownSkillLevels() string {
	ss := []string{}
	for _, need := range needs.AsArray() {
		ss = append(ss, fmt.Sprintf("%s", need.SkillLevel()))
	}
	return strings.Join(ss, ", ")
}

func (needs Needs) AsArray() []Need {
	a := []Need{}
	for _, skillLevel := range needs.IDs() {
		a = append(a, needs.Get(skillLevel))
	}
	return a
}

func (needs Needs) Unmet(users *Users) *Needs {
	out := &needs
	for _, user := range users.AsArray() {
		out = out.CheckRequired(user)
	}
	return out
}

func (needs *Needs) CheckLimits(user *User) (adjusted, modified, violated *Needs) {
	violated = NewNeeds()
	modified = NewNeeds()
	adjusted = NewNeeds()
	for _, need := range needs.AsArray() {
		qualified, adjustedNeed := need.QualifyUser(user)
		if !qualified {
			adjusted.Set(need)
			continue
		}
		adjusted.Set(adjustedNeed)
		modified.Set(adjustedNeed)
		if adjustedNeed.Count() < 0 {
			violated.Set(need)
		}
	}

	return adjusted, modified, violated
}

func (require *Needs) CheckRequired(user *User) (adjusted *Needs) {
	adjusted = NewNeeds()
	for _, need := range require.AsArray() {
		qualified, adjustedNeed := need.QualifyUser(user)
		if !qualified {
			adjusted.Set(need)
			continue
		}
		adjusted.Set(adjustedNeed)
	}

	return adjusted
}
