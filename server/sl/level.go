// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Level int

type SkillLevel struct {
	Skill string
	Level Level
}

var _ pflag.Value = (*Level)(nil)

const (
	Beginner = Level(1) + iota
	Intermediate
	Advanced
	Expert
)

const (
	LevelAdvanced           = "advanced"
	LevelAdvancedSymbol     = "◈"
	LevelBeginner           = "beginner"
	LevelBeginnerSymbol     = "◉"
	LevelExpert             = "expert"
	LevelExpertSymbol       = "◈◈"
	LevelIntermediate       = "intermediate"
	LevelIntermediateSymbol = "▣"
)

func (level Level) String() string {
	switch level {
	case Beginner:
		return LevelBeginnerSymbol
	case Intermediate:
		return LevelIntermediateSymbol
	case Advanced:
		return LevelAdvancedSymbol
	case Expert:
		return LevelExpertSymbol
	}
	return "⛔"
}

func (level *Level) Type() string {
	return "skill_level"
}

func (level *Level) Set(in string) error {
	switch in {
	case LevelBeginner,
		LevelBeginnerSymbol,
		LevelBeginner + LevelBeginnerSymbol,
		"1":
		*level = Beginner

	case LevelIntermediate,
		LevelIntermediateSymbol,
		LevelIntermediate + LevelIntermediateSymbol,
		"2":
		*level = Intermediate

	case LevelAdvanced,
		LevelAdvancedSymbol,
		LevelAdvanced + LevelAdvancedSymbol,
		"3":
		*level = Advanced

	case LevelExpert,
		LevelExpertSymbol,
		LevelExpert + LevelExpertSymbol,
		"4":
		*level = Expert

	default:
		return errors.Errorf("%s is not a valid skill level", in)
	}
	return nil
}

func (skillLevel SkillLevel) String() string {
	return fmt.Sprintf("%s-%s", skillLevel.Skill, skillLevel.Level)
}

func (skillLevel SkillLevel) Type() string {
	return fmt.Sprintf("%s-%v", skillLevel.Skill, int(skillLevel.Level))
}

func (skillLevel *SkillLevel) Set(in string) error {
	split := strings.Split(in, "-")
	if len(split) > 2 {
		return errors.Errorf("should be formatted as skill-level: %s", skillLevel)
	}

	l := Beginner
	if len(split) == 2 {
		err := l.Set(split[1])
		if err != nil {
			return err
		}
	}
	skillLevel.Skill = split[0]
	skillLevel.Level = l
	return nil
}
