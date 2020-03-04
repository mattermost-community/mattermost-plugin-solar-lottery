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
	*types.IntSet
}

type Need struct {
	SkillLevel
	Count int
}

func NewNeeds() Needs {
	return Needs{
		IntSet: types.NewIntSet(),
	}
}

func (needs Needs) Markdown() string {
	ss := []string{}
	for _, skillLevel := range needs.IDs() {
		ss = append(ss, fmt.Sprintf("**%v** %s", needs.Get(skillLevel), skillLevel))
	}
	return strings.Join(ss, ", ")
}

func (needs Needs) UnmetRequirements(users *Users) Needs {
	work := NewNeeds()
	for _, id := range needs.IDs() {
		work.Set(id, needs.Get(id))
	}

	for _, id := range work.IDs() {
		skillLevel := ParseSkillLevel(id)
		for _, user := range users.AsArray() {
			if user.IsQualified(skillLevel) {
				work.Set(id, work.Get(id)-1)
			}
		}
	}

	unmet := NewNeeds()
	for _, key := range work.IDs() {
		count := work.Get(key)
		if count > 0 {
			unmet.Set(key, count)
		}
	}

	return unmet
}

func NewNeed(count int, skillLevel SkillLevel) *Need {
	return &Need{
		Count:      count,
		SkillLevel: skillLevel,
	}
}

func (need Need) GetID() types.ID {
	return types.ID(need.SkillLevel.String())
}

func (need Need) String() string {
	return fmt.Sprintf("%v-%s", need.Count, need.SkillLevel)
}

func (need Need) Markdown() string {
	return fmt.Sprintf("**%v** %s", need.Count, need.SkillLevel)
}
