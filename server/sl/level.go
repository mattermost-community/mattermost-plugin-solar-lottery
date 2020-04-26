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
	AnyLevel = Level(iota)
	BeginnerLevel
	IntermediateLevel
	AdvancedLevel
	ExpertLevel
)

const (
	AdvancedLevelName       = "advanced"
	AdvancedLevelSymbol     = "◈"
	AnyLevelName            = "any"
	AnyLevelSymbol          = "*"
	BeginnerLevelName       = "beginner"
	BeginnerLevelSymbol     = "◉"
	ExpertLevelName         = "expert"
	ExpertLevelSymbol       = "◈◈"
	IntermediateLevelName   = "intermediate"
	IntermediateLevelSymbol = "▣"
)

func (level Level) String() string {
	switch level {
	case AnyLevel:
		return AnyLevelSymbol
	case BeginnerLevel:
		return BeginnerLevelSymbol
	case IntermediateLevel:
		return IntermediateLevelSymbol
	case AdvancedLevel:
		return AdvancedLevelSymbol
	case ExpertLevel:
		return ExpertLevelSymbol
	}
	return "⛔"
}

func (level *Level) Type() string {
	return "level"
}

func (level *Level) Set(in string) error {
	switch in {
	case BeginnerLevelName, BeginnerLevelSymbol, "1":
		*level = BeginnerLevel
	case IntermediateLevelName, IntermediateLevelSymbol, "2":
		*level = IntermediateLevel
	case AdvancedLevelName, AdvancedLevelSymbol, "3":
		*level = AdvancedLevel
	case ExpertLevelName, ExpertLevelSymbol, "4":
		*level = ExpertLevel
	default:
		return errors.Errorf("%s is not a valid skill level", in)
	}
	return nil
}
