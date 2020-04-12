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

func TestCommandTaskNewTicket(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param ticket test-rotation
			/lotto rotation param min -s webapp-2 --count 2 test-rotation
			/lotto rotation param max -s server-3 --count 1 test-rotation
			`)
		require.NoError(t, err)

		out, err := runTaskCreateCommand(t, SL, `/lotto task new ticket test-rotation --summary test-summary1`)
		task := out.Task
		require.NoError(t, err)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary1", task.Summary)
		require.Equal(t, "0001-01-01", task.ExpectedStart.String())

		require.Equal(t, map[types.ID]int64{"*-*": 1, "webapp-▣": 2}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◈": 1}, task.Limit.TestAsMap())
	})
}
