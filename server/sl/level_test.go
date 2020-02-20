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
		{Level(0), "⛔"},
		{Level(1), LevelBeginnerSymbol},
		{Level(2), LevelIntermediateSymbol},
		{Level(3), LevelAdvancedSymbol},
		{Level(4), LevelExpertSymbol},
	} {
		assert.Equal(t, tc.expected, fmt.Sprint(tc.level))
	}

	for _, tc := range []struct {
		value    string
		expected Level
	}{
		// {LevelBeginnerSymbol, Level(1)},
		{LevelIntermediateSymbol, Level(2)},
		// {LevelAdvancedSymbol, Level(3)},
		// {LevelExpertSymbol, Level(4)},
	} {
		var l Level
		err := l.Set(tc.value)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, l)
	}

	for _, bad := range []string{"⛔", "invalid"} {
		l := Beginner
		err := l.Set(bad)
		require.Error(t, err)
		require.Equal(t, Beginner, l)
	}
}
