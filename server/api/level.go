// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/pkg/errors"

const (
	Beginner = 1 + iota
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

func ParseLevel(l string) (int, error) {
	switch l {
	case LevelBeginner, LevelBeginnerSymbol:
		return Beginner, nil
	case LevelIntermediate, LevelIntermediateSymbol:
		return Intermediate, nil
	case LevelAdvanced, LevelAdvancedSymbol:
		return Advanced, nil
	case LevelExpert, LevelExpertSymbol:
		return Expert, nil
	}
	return 0, errors.Errorf("%s is not a valid skill level", l)
}

func LevelToString(l int) string {
	switch l {
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
