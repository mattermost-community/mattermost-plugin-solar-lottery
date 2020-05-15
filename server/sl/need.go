// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Need struct {
	types.IntValue
}

func NewNeed(count int64, skillLevel SkillLevel) Need {
	return Need{
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
	return fmt.Sprintf("%v %s", need.Count(), need.SkillLevel())
}

func (need Need) Markdown() md.MD {
	return md.Markdownf("**%v** %s", need.Count(), need.SkillLevel())
}

func (need Need) QualifyUser(user *User) (bool, Need) {
	skillLevel := need.SkillLevel()
	skill := skillLevel.Skill
	level := int64(skillLevel.Level)
	ulevel := user.SkillLevels.Get(skill)

	if ulevel >= level {
		need.Value--
		return true, need
	}
	return false, need
}

func (need Need) QualifyUsers(users *Users) (*Users, Need) {
	qualified := NewUsers()
	for _, user := range users.AsArray() {
		isQualified, adj := need.QualifyUser(user)
		if !isQualified {
			continue
		}
		need = adj
		qualified.Set(user)
	}
	return qualified, need
}
