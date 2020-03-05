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

func NewNeeds(nn ...*Need) *Needs {
	needs := &Needs{
		IntSet: *types.NewIntSet(),
	}
	for _, need := range nn {
		needs.Set(need)
	}
	return needs
}

func (needs Needs) Get(id types.ID) *Need {
	if !needs.Contains(id) {
		return nil
	}
	return NewNeed(needs.IntSet.Get(id), ParseSkillLevel(id))
}

func (needs Needs) GetCountForSkillLevel(skillLevel SkillLevel) int64 {
	id := types.ID(skillLevel.String())
	return needs.IntSet.Get(id)
}

func (needs *Needs) Set(need *Need) {
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

func (needs Needs) AsArray() []*Need {
	a := []*Need{}
	for _, skillLevel := range needs.IDs() {
		a = append(a, needs.Get(skillLevel))
	}
	return a
}

type Need struct {
	types.IntValue
}

func NewNeed(count int64, skillLevel SkillLevel) *Need {
	return &Need{
		IntValue: types.NewIntValue(types.ID(skillLevel.String()), count),
	}
}

func (need Need) Count() int64 {
	return need.Value
}

func (need Need) SkillLevel() SkillLevel {
	return ParseSkillLevel(need.GetID())
}

func (need Need) String() string {
	return fmt.Sprintf("%v-%s", need.Count(), need.SkillLevel())
}

func (need Need) Markdown() string {
	return fmt.Sprintf("**%v** %s", need.Count(), need.SkillLevel())
}

func (need Need) QualifyUser(user *User) (bool, *Need) {
	adjusted := need
	skillLevel := need.SkillLevel()
	skill := skillLevel.Skill
	level := int64(skillLevel.Level)
	ulevel := user.SkillLevels.Get(skill)

	if ulevel >= level {
		adjusted.Value--
		return true, &adjusted
	}
	return false, nil
}

func (need Need) QualifyUsers(users *Users) (*Users, *Need) {
	adjusted := &need
	qualified := NewUsers()

	for _, user := range users.AsArray() {
		isQualified, adj := adjusted.QualifyUser(user)
		if !isQualified {
			continue
		}
		adjusted = adj
		qualified.Set(user)
	}

	return qualified, adjusted
}
