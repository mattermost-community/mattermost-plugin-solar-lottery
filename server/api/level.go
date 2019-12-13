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
	LevelBeginner     = "beginner"
	LevelIntermediate = "intermediate"
	LevelAdvanced     = "advanced"
	LevelExpert       = "expert"
)

func ParseLevel(l string) (int, error) {
	switch l {
	case LevelBeginner:
		return Beginner, nil
	case LevelIntermediate:
		return Intermediate, nil
	case LevelAdvanced:
		return Advanced, nil
	case LevelExpert:
		return Expert, nil
	}
	return 0, errors.Errorf("%s is not a valid skill level", l)
}

func LevelToString(l int) string {
	switch l {
	case Beginner:
		return LevelBeginner
	case Intermediate:
		return LevelIntermediate
	case Advanced:
		return LevelAdvanced
	case Expert:
		return LevelExpert
	}
	return "none"
}
