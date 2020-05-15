// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestUserDisqualify(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto user qualify -s webapp-2
			/lotto user qualify -s somethingelse-3
		`)

		users := mustRunUsersQualify(t, SL, `/lotto user disqualify -s webapp`)
		require.Equal(t, 1, len(users.TestArray()))
		u := users.TestArray()[0]
		require.Equal(t, types.ID("test-user"), u.MattermostUserID)
		require.EqualValues(t, types.NewIntSet(types.NewIntValue("somethingelse", 3)), u.SkillLevels)
	})
}
