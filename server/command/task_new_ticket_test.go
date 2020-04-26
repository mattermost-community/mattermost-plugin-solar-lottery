// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskNewTicket(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=ticket
			/lotto rotation set require -s webapp-2 --count 2 test-rotation
			/lotto rotation set limit -s server-3 --count 1 test-rotation
			`)

		task := mustRunTaskCreate(t, SL, `/lotto task new ticket test-rotation --summary test-summary1 --now=2020-03-03`)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary1", task.Summary)
		require.Equal(t, "2020-03-03T08:00", task.ExpectedStart.String())

		require.Equal(t, map[types.ID]int64{"any": 1, "webapp-▣": 2}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◈": 1}, task.Limit.TestAsMap())
	})
}
