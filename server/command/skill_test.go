// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/stretchr/testify/require"
)

func TestSkillList(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)
		mustRunMulti(t, SL, `
			/lotto skill new test
			/lotto skill new test-123
			/lotto skill new test-345
			`)

		knownSkills, err := store.IDIndex(sl.KeyKnownSkills).Load()
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-123", "test-345"}, knownSkills.TestIDs())

		out := []string{}
		mustRunJSON(t, SL, `/lotto skill list`, &out)
		require.Equal(t, []string{"test", "test-123", "test-345"}, out)
	})
}
