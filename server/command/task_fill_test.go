// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskFill(t *testing.T) {
	t.Run("fill happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param ticket test-rotation 
			/lotto rotation param min -s web-1 --count 2 test-rotation
			
			# user3 and 4 are joining in the future, and will not be selected,
			# user1 and 2 are in the past, and will be selected
			/lotto user join test-rotation @test-user3 @test-user4 --start 2033-01-01
			/lotto user join test-rotation @test-user1 @test-user2 --start 2020-01-01
			/lotto user qualify -s web-1 @test-user1 @test-user2 @test-user3 @test-user4
			/lotto task new ticket test-rotation --summary test-summary1
			/lotto task show test-rotation#1
			`)
		require.NoError(t, err)

		out := &sl.OutAssignTask{
			Changed: sl.NewUsers(),
		}
		_, err = runJSONCommand(t, SL, `
			/lotto task fill test-rotation#1 --debug-now 2020-03-03
			`, &out)
		task := out.Task
		require.NoError(t, err)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary1", task.Summary)
		require.Equal(t, []string{"test-user1", "test-user2"}, out.Task.MattermostUserIDs.TestIDs())

		user := sl.NewUser("")
		_, err = runJSONCommand(t, SL, `
			/lotto user show @test-user1
		`, &user)
		require.NoError(t, err)
		// compare in UTC, just in case
		require.Equal(t, types.MustParseTime("2020-03-03T08:00").Time, time.Unix(user.LastServed.Get("test-rotation"), 0).UTC())
	})
}
