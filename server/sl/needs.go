// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strconv"
)

type Need struct {
	Count int
	Skill string
	Level int
}

func NewNeed(skill string, level int, count int) *Need {
	return &Need{
		Count: count,
		Skill: skill,
		Level: level,
	}
}

func (need Need) String() string {
	return fmt.Sprintf("%s-%v-%v", need.Skill, need.Level, need.Count)
}

func (need Need) SkillLevel() string {
	return need.Skill + "-" + strconv.Itoa(need.Level)
}

func (need *Need) Markdown() string {
	return fmt.Sprintf("**%v** %s", need.Count, need.SkillLevel())
}

func UnmetRequirements(needs []*Need, users UserMap) []*Need {
	work := append([]*Need{}, needs...)
	for i, need := range work {
		for _, user := range users {
			if user.IsQualified(need) {
				work[i].Count--
			}
		}
	}

	var unmet []*Need
	for _, need := range work {
		if need.Count > 0 {
			unmet = append(unmet, need)
		}
	}

	return unmet
}
