// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

type SkillLevel struct {
	Skill types.ID
	Level Level
}

var AnySkill = types.ID("any")
var AnySkillLevel = NewSkillLevel(AnySkill, AnyLevel)

func NewSkillLevel(s types.ID, l Level) SkillLevel {
	return SkillLevel{
		Skill: s,
		Level: l,
	}
}

func ParseSkillLevel(in types.ID) SkillLevel {
	skillLevel := SkillLevel{}
	_ = skillLevel.Set(string(in))
	return skillLevel
}

func (skillLevel SkillLevel) String() string {
	none := SkillLevel{}
	switch {
	case skillLevel == none || skillLevel.Skill == AnySkill:
		return AnySkill.String()

	case skillLevel.Level == AnyLevel:
		return skillLevel.Skill.String()

	default:
		return fmt.Sprintf("%s-%s", skillLevel.Skill, skillLevel.Level)
	}
}

func (skillLevel SkillLevel) AsID() types.ID {
	return types.ID(skillLevel.String())
}

func (skillLevel SkillLevel) Type() string {
	return fmt.Sprintf("%s-%v", skillLevel.Skill, int(skillLevel.Level))
}

func (skillLevel *SkillLevel) Set(in string) error {
	if in == AnySkill.String() {
		*skillLevel = AnySkillLevel
		return nil
	}
	split := strings.Split(in, "-")
	if len(split) > 2 {
		return errors.Errorf("should be formatted as skill-level: %s", skillLevel)
	}

	l := BeginnerLevel
	if len(split) == 2 {
		err := l.Set(split[1])
		if err != nil {
			return err
		}
	}
	skillLevel.Skill = types.ID(split[0])
	skillLevel.Level = l
	return nil
}

func (skillLevel SkillLevel) GetID() types.ID {
	return skillLevel.Skill
}
