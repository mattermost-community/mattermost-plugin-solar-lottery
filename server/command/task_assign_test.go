// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskAssign(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()

		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=ticket
			/lotto rotation set require -s webapp-2 --count 2 test-rotation
			/lotto rotation set limit -s server-3 --count 1 test-rotation
			/lotto task new ticket test-rotation --summary test-summary1
			/lotto task new ticket test-rotation --summary test-summary2
			`)

		task := mustRunTask(t, SL, `/lotto task show test-rotation#2`)
		require.Equal(t, map[types.ID]int64{"any": 1, "webapp-▣": 2}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◈": 1}, task.Limit.TestAsMap())

		task = mustRunTaskAssign(t, SL, `/lotto task assign test-rotation#2 @test-user1 @test-user2`)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#2"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary2", task.Summary)

		// No change in needs
		require.Equal(t, []string{"any", "webapp-▣"}, task.Require.TestIDs())
		require.Equal(t, int64(2), task.Require.IntSet.Get("webapp-▣"))
		require.Equal(t, []string{"server-◈"}, task.Limit.TestIDs())
		require.Equal(t, int64(1), task.Limit.IntSet.Get("server-◈"))

		require.Equal(t, []string{"test-user1", "test-user2"}, task.MattermostUserIDs.TestIDs())
	})

	t.Run("max constraint", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()

		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=ticket
			/lotto rotation set limit -s server-1 --count 1 test-rotation
			/lotto user qualify @test-user1 @test-user2 -s server-1
			/lotto task new ticket test-rotation --summary test-summary1
			`)

		_, err := run(t, SL, `/lotto task assign test-rotation#1 @test-user1 @test-user2`)
		require.Error(t, err)
		require.Equal(t, "failed to assign task test-rotation#1: user @test-user2 failed max constraints server-◉", err.Error())
	})

	t.Run("max constraint--force", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()

		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=ticket
			/lotto rotation set limit -s server-1 --count 1 test-rotation
			/lotto user qualify @test-user1 @test-user2 -s server-1
			/lotto task new ticket test-rotation --summary test-summary1
			`)

		task := mustRunTaskAssign(t, SL, `/lotto task assign test-rotation#1 @test-user1 @test-user2 --force`)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary1", task.Summary)
		require.Equal(t, []string{"server-◉"}, task.Limit.TestIDs())
		require.Equal(t, int64(1), task.Limit.IntSet.Get("server-◉"))
		require.Equal(t, []string{"test-user1", "test-user2"}, task.MattermostUserIDs.TestIDs())
	})
}
