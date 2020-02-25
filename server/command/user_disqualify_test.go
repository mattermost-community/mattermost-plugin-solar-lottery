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

func TestCommandUserDisqualify(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto user qualify -k webapp-2
			/lotto user qualify -k somethingelse-3
		`)
		require.NoError(t, err)

		uu := sl.UserMap{}
		_, err = runJSONCommand(t, SL, `
			/lotto user disqualify -k webapp`, &uu)
		require.NoError(t, err)
		require.EqualValues(t, []sl.User{
			sl.User{
				PluginVersion:    "test-plugin-version",
				MattermostUserID: "test-user",
				SkillLevels: types.IntMap{
					"somethingelse": 3,
				},
			},
		}, uu.Sorted())
	})
}
