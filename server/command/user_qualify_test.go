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

		qualified := sl.NewUsers()
		_, err := runJSONCommand(t, SL, `
			/lotto user qualify @uid1-username -k webapp-▣ @uid2-username`, &qualified)
		require.NoError(t, err)
		require.Equal(t, []string{"uid1", "uid2"}, qualified.TestIDs())

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

		qualified = sl.NewUsers()
		_, err = runJSONCommand(t, SL, `
			/lotto user qualify -k somethingelse-3`, &qualified)
		require.NoError(t, err)
		require.Equal(t, []string{"test-user"}, qualified.TestIDs())

	})
}
