// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevelHappy(t *testing.T) {
	for _, tc := range []struct {
		level    Level
		expected string
	}{
		{Level(-1), "⛔"},
		{Level(0), AnyLevelSymbol},
		{Level(1), BeginnerLevelSymbol},
		{Level(2), IntermediateLevelSymbol},
		{Level(3), AdvancedLevelSymbol},
		{Level(4), ExpertLevelSymbol},
	} {
		assert.Equal(t, tc.expected, fmt.Sprint(tc.level))
	}

	for _, tc := range []struct {
		value    string
		expected Level
	}{
		{BeginnerLevelSymbol, Level(1)},
		{IntermediateLevelSymbol, Level(2)},
		{AdvancedLevelSymbol, Level(3)},
		{ExpertLevelSymbol, Level(4)},
	} {
		var l Level
		err := l.Set(tc.value)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, l)
	}

	for _, bad := range []string{"⛔", "invalid"} {
		l := BeginnerLevel
		err := l.Set(bad)
		require.Error(t, err)
		require.Equal(t, BeginnerLevel, l)
	}
}
