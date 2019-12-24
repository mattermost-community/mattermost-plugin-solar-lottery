// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevelHappy(t *testing.T) {
	var l Level
	err := l.Set(LevelExpert)
	require.NoError(t, err)
	assert.Equal(t, LevelExpertSymbol, fmt.Sprint(l))
}
