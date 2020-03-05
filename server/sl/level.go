// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Level int64

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
