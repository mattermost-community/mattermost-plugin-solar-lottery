// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
)

func TestCommandUserQualify(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		qualified := sl.UserMap{}
		_, err := runJSONCommand(t, SL, `
			/lotto user qualify @uid1-username -k webapp -l 2 @uid2-username`, &qualified)
		require.NoError(t, err)
		require.Equal(t, []string{"uid1", "uid2"}, qualified.IDs().Sorted())

		uu := sl.UserMap{}
		_, err = runJSONCommand(t, SL, `
			/lotto user show @uid1-username @uid2-username`, &uu)
		require.NoError(t, err)
		require.EqualValues(t, []sl.User{
			sl.User{
				PluginVersion:    "test-plugin-version",
				MattermostUserID: "uid1",
				SkillLevels: sl.IntMap{
					"webapp": 2,
				},
			},
			sl.User{
				PluginVersion:    "test-plugin-version",
				MattermostUserID: "uid2",
				SkillLevels: sl.IntMap{
					"webapp": 2,
				},
			},
		}, uu.Sorted())

		qualified = sl.UserMap{}
		_, err = runJSONCommand(t, SL, `
			/lotto user qualify -k somethingelse -l 2`, &qualified)
		require.NoError(t, err)
		require.Equal(t, []string{"test-user"}, qualified.IDs().Sorted())

	})
}
