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

		qualified := sl.UserMap{}
		_, err := runJSONCommand(t, SL, `
			/lotto user qualify @uid1-username -k webapp-▣ @uid2-username`, &qualified)
		require.NoError(t, err)
		require.Equal(t, []string{"uid1", "uid2"}, qualified.IDs().TestIDs())

		uu := sl.UserMap{}
		_, err = runJSONCommand(t, SL, `
			/lotto user show @uid1-username @uid2-username`, &uu)
		require.NoError(t, err)
		require.Equal(t, 2, len(uu.TestSorted()))
		u1 := uu.TestSorted()[0]
		require.Equal(t, types.ID("uid1"), u1.MattermostUserID)
		require.EqualValues(t, types.NewIntIndex(types.NewIDInt("webapp", 2)), u1.SkillLevels)

		u2 := uu.TestSorted()[1]
		require.Equal(t, types.ID("uid2"), u2.MattermostUserID)
		require.EqualValues(t, types.NewIntIndex(types.NewIDInt("webapp", 2)), u2.SkillLevels)

		qualified = sl.UserMap{}
		_, err = runJSONCommand(t, SL, `
			/lotto user qualify -k somethingelse-3`, &qualified)
		require.NoError(t, err)
		require.Equal(t, []string{"test-user"}, qualified.IDs().TestIDs())

	})
}
