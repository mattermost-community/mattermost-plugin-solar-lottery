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

func TestCommandTaskAssign(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param ticket test-rotation
			/lotto rotation param min -s webapp-2 --count 2 test-rotation
			/lotto rotation param max -s server-3 --count 1 test-rotation
			/lotto task new ticket test-rotation --summary test-summary1
			/lotto task new ticket test-rotation --summary test-summary2
			`)
		require.NoError(t, err)

		task := sl.NewTask("")
		_, err = runJSONCommand(t, SL, `
			/lotto task show test-rotation#2
			`, &task)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 1, "webapp-▣": 2}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◈": 1}, task.Limit.TestAsMap())

		out := &sl.OutAssignTask{
			Changed: sl.NewUsers(),
		}
		_, err = runJSONCommand(t, SL, `
			/lotto task assign test-rotation#2 @test-user1 @test-user2
			`, &out)
		task = out.Task
		require.NoError(t, err)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#2"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary2", task.Summary)

		// No change in needs
		require.Equal(t, []string{"*-*", "webapp-▣"}, task.Require.TestIDs())
		require.Equal(t, int64(2), task.Require.IntSet.Get("webapp-▣"))
		require.Equal(t, []string{"server-◈"}, task.Limit.TestIDs())
		require.Equal(t, int64(1), task.Limit.IntSet.Get("server-◈"))

		require.Equal(t, []string{"test-user1", "test-user2"}, task.MattermostUserIDs.TestIDs())
	})

	t.Run("max constraint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param ticket test-rotation
			/lotto rotation param max -s server-1 --count 1 test-rotation
			/lotto user qualify @test-user1 @test-user2 -s server-1
			/lotto task new ticket test-rotation --summary test-summary1
			`)
		require.NoError(t, err)

		_, err = runJSONCommand(t, SL, `
			/lotto task assign test-rotation#1 @test-user1 @test-user2
			`, nil)
		require.Error(t, err)
		require.Equal(t, "user @test-user2-username failed max constraints server-◉", err.Error())
	})

	t.Run("max constraint--force", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param ticket test-rotation
			/lotto rotation param max -s server-1 --count 1 test-rotation
			/lotto user qualify @test-user1 @test-user2 -s server-1
			/lotto task new ticket test-rotation --summary test-summary1
			`)
		require.NoError(t, err)

		out := &sl.OutAssignTask{
			Changed: sl.NewUsers(),
		}
		_, err = runJSONCommand(t, SL, `
			/lotto task assign test-rotation#1 @test-user1 @test-user2 --force
			`, &out)
		require.NoError(t, err)
		require.Equal(t, "test-plugin-version", out.Task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#1"), out.Task.TaskID)
		require.Equal(t, types.ID("test-rotation"), out.Task.RotationID)
		require.Equal(t, sl.TaskStatePending, out.Task.State)
		require.Equal(t, "test-summary1", out.Task.Summary)
		require.Equal(t, []string{"server-◉"}, out.Task.Limit.TestIDs())
		require.Equal(t, int64(1), out.Task.Limit.IntSet.Get("server-◉"))
		require.Equal(t, []string{"test-user1", "test-user2"}, out.Task.MattermostUserIDs.TestIDs())
	})
}
