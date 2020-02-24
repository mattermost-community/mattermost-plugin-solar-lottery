// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/stretchr/testify/require"
)

func TestCommandSkillList(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto skill add test
			/lotto skill add test-123
			/lotto skill add test-345
			`)
		require.NoError(t, err)

		knownSkills, err := store.Index(sl.KeyKnownSkills).Load()
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-123", "test-345"}, knownSkills.Sorted())

		out := []string{}
		_, err = runJSONCommand(t, SL, `
			/lotto skill list`, &out)
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-123", "test-345"}, out)

		// /lotto skill delete test
	})
}
