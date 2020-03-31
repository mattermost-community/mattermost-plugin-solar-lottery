// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestCommandUserQualify(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		out := sl.OutQualify{
			Users: sl.NewUsers(),
		}
		_, err := runJSONCommand(t, SL, `
			/lotto user qualify @uid1-username -s webapp-▣ @uid2-username`, &out)
		require.NoError(t, err)
		require.Equal(t, []string{"uid1", "uid2"}, out.Users.TestIDs())

		ss := []string{}
		_, err = runJSONCommand(t, SL, `
			/lotto skill list`, &ss)
		require.NoError(t, err)
		require.Equal(t, []string{"webapp"}, ss)

		uu := sl.NewUsers()
		_, err = runJSONCommand(t, SL, `
			/lotto user show @uid1-username @uid2-username`, &uu)
		require.NoError(t, err)
		require.Equal(t, 2, len(uu.TestIDs()))
		u1 := uu.TestArray()[0]
		require.Equal(t, types.ID("uid1"), u1.MattermostUserID)
		require.EqualValues(t, types.NewIntSet(types.NewIntValue("webapp", 2)), u1.SkillLevels)

		u2 := uu.TestArray()[1]
		require.Equal(t, types.ID("uid2"), u2.MattermostUserID)
		require.EqualValues(t, types.NewIntSet(types.NewIntValue("webapp", 2)), u2.SkillLevels)

		out = sl.OutQualify{
			Users: sl.NewUsers(),
		}
		_, err = runJSONCommand(t, SL, `
			/lotto user qualify -s somethingelse-3`, &out)
		require.NoError(t, err)
		require.Equal(t, []string{"test-user"}, out.Users.TestIDs())
		require.Equal(t, int64(3), out.Users.Get("test-user").SkillLevels.Get("somethingelse"))
	})
}
