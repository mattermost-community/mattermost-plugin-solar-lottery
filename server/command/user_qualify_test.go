// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestUserQualify(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()

		users := mustRunUsersQualify(t, SL,
			`/lotto user qualify @uid1 -s webapp-▣ @uid2`)
		require.Equal(t, []string{"uid1", "uid2"}, users.TestIDs())

		ss := []string{}
		mustRunJSON(t, SL,
			`/lotto skill list`, &ss)
		require.Equal(t, []string{"webapp"}, ss)

		users = mustRunUsers(t, SL,
			`/lotto user show @uid1 @uid2`)
		require.Equal(t, 2, len(users.TestIDs()))
		u1 := users.TestArray()[0]
		require.Equal(t, types.ID("uid1"), u1.MattermostUserID)
		require.EqualValues(t, types.NewIntSet(types.NewIntValue("webapp", 2)), u1.SkillLevels)
		u2 := users.TestArray()[1]
		require.Equal(t, types.ID("uid2"), u2.MattermostUserID)
		require.EqualValues(t, types.NewIntSet(types.NewIntValue("webapp", 2)), u2.SkillLevels)

		users = mustRunUsersQualify(t, SL,
			`/lotto user qualify -s somethingelse-3`)
		require.Equal(t, []string{"test-user"}, users.TestIDs())
		require.Equal(t, int64(3), users.Get("test-user").SkillLevels.Get("somethingelse"))
	})
}
