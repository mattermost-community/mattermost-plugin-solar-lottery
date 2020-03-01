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

		uu := sl.NewUsers()
		_, err = runJSONCommand(t, SL, `
			/lotto user disqualify -k webapp`, &uu)
		require.NoError(t, err)
		require.Equal(t, 1, len(uu.TestArray()))
		u := uu.TestArray()[0]
		require.Equal(t, types.ID("test-user"), u.MattermostUserID)
		require.EqualValues(t, types.NewIntIndex(types.NewIDInt("somethingelse", 3)), u.SkillLevels)
	})
}
