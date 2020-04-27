// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskFill(t *testing.T) {
	t.Run("fill happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=ticket --beginning=2020-03-01
			/lotto rotation set require -s web-1 --count 2 test-rotation
			
			# user1,2 are joining in the future, and will not be selected,
			# user3,4 are in the past, and will be selected

			/lotto user join test-rotation @test-user1 --starting 2033-01-01
			/lotto user join test-rotation @test-user2 --starting 2033-01-01PST
			/lotto user join test-rotation @test-user3 --starting 2020-01-01UTC
			/lotto user join test-rotation @test-user4 --starting 2020-01-01T11:00EST
			/lotto user qualify -s web-1 @test-user1 @test-user2 @test-user3 @test-user4
			/lotto task new ticket test-rotation --summary test-summary1 --now 2020-03-01
			`)

		lastServed := func(user *sl.User) string {
			return time.Unix(user.LastServed.Get("test-rotation"), 0).Format(time.RFC3339)
		}
		require.Equal(t, "2033-01-01T00:00:00-08:00", lastServed(mustRunUser(t, SL, `/lotto user show @test-user1`)))
		require.Equal(t, "2033-01-01T00:00:00-08:00", lastServed(mustRunUser(t, SL, `/lotto user show @test-user2`)))
		require.Equal(t, "2019-12-31T16:00:00-08:00", lastServed(mustRunUser(t, SL, `/lotto user show @test-user3`)))
		require.Equal(t, "2020-01-01T03:00:00-08:00", lastServed(mustRunUser(t, SL, `/lotto user show @test-user4`)))

		task := mustRunTaskAssign(t, SL, `/lotto task fill test-rotation#1 --now 2020-02-20`)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "test-summary1", task.Summary)
		require.Equal(t, []string{"test-user3", "test-user4"}, task.MattermostUserIDs.TestIDs())
		require.Equal(t, "2020-03-01T08:00", task.ExpectedStart.String())
		require.Equal(t, "30m0s", task.ExpectedDuration.String())

		// Once added to a ticket, last served must be updated
		require.Equal(t, "2020-03-01T00:30:00-08:00", lastServed(mustRunUser(t, SL, `/lotto user show @test-user3`)))
		require.Equal(t, "2020-03-01T00:30:00-08:00", lastServed(mustRunUser(t, SL, `/lotto user show @test-user4`)))

		// Check at least one calendar
		require.Equal(t,
			[]*sl.Unavailable{
				{
					Interval: types.Interval{
						Start:  types.MustParseTime("2020-03-01T08:00"),
						Finish: types.MustParseTime("2020-03-01T08:30"),
					},
					Reason:     "task",
					TaskID:     "test-rotation#1",
					RotationID: "test-rotation",
				},
			},
			mustRunUser(t, SL, `/lotto user show @test-user3`).Calendar,
		)
	})
}
