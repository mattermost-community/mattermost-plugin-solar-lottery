// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Level int

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
	return "none"
}

func (level *Level) Type() string {
	return "skill_level"
}

func (level *Level) Set(in string) error {
	switch in {
	case LevelBeginner,
		LevelBeginnerSymbol,
		LevelBeginner + LevelBeginnerSymbol:
		*level = Beginner

	case LevelIntermediate,
		LevelIntermediateSymbol,
		LevelIntermediate + LevelIntermediateSymbol:
		*level = Intermediate

	case LevelAdvanced,
		LevelAdvancedSymbol,
		LevelAdvanced + LevelAdvancedSymbol:
		*level = Advanced

	case LevelExpert,
		LevelExpertSymbol,
		LevelExpert + LevelExpertSymbol:
		*level = Expert

	default:
		return errors.Errorf("%s is not a valid skill level", in)
	}
	return nil
}
