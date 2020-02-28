// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Need struct {
	SkillLevel
	Count int
}

func (need Need) GetID() string {
	return need.SkillLevel.String()
}

func (need Need) Clone(bool) types.Cloneable {
	return need
}

type Needs struct {
	*types.IntIndex
}

type NeedSetProto []*Need

func NewNeeds(needs ...*Need) Needs {
	return Needs{
		IntIndex: types.NewIntIndex(),
	}
}

func (needs Needs) Markdown() string {
	ss := []string{}
	for _, skillLevel := range needs.IDs() {
		ss = append(ss, fmt.Sprintf("**%v** %s", needs.Get(skillLevel), skillLevel))
	}
	return strings.Join(ss, ", ")
}

func (needs Needs) UnmetRequirements(users UserMap) Needs {
	work := needs.Clone(false).(Needs)
	for _, id := range work.IDs() {
		skillLevel := SkillLevel{}
		_ = skillLevel.Set(string(id))
		for _, user := range users {
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

func (need Need) String() string {
	return fmt.Sprintf("%v-%s", need.Count, need.SkillLevel)
}

func (need Need) Markdown() string {
	return fmt.Sprintf("**%v** %s", need.Count, need.SkillLevel)
}
