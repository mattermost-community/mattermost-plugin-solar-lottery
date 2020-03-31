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
			/lotto user qualify -s webapp-2
			/lotto user qualify -s somethingelse-3
		`)
		require.NoError(t, err)

		out := sl.OutQualify{
			Users: sl.NewUsers(),
		}
		_, err = runJSONCommand(t, SL, `
			/lotto user disqualify -s webapp`, &out)
		require.NoError(t, err)
		require.Equal(t, 1, len(out.Users.TestArray()))
		u := out.Users.TestArray()[0]
		require.Equal(t, types.ID("test-user"), u.MattermostUserID)
		require.EqualValues(t, types.NewIntSet(types.NewIntValue("somethingelse", 3)), u.SkillLevels)
	})
}
