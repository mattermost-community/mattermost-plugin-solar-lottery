// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskTransition(t *testing.T) {
	t.Run("user messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		poster := &bot.TestPoster{}
		SL, _ := getTestSLWithPoster(t, ctrl, poster)
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=ticket
			/lotto user join test-rotation @test-user1 @test-user2 @test-user3 @test-user4 @test-user5 @test-user6
			`)
		poster.Reset()

		mustRunMulti(t, SL, `
			/lotto task new ticket test-rotation --summary test-summary1
			/lotto task assign test-rotation#1 @test-user3 @test-user5
			/lotto task schedule test-rotation#1 
			`)
		require.Equal(t, []bot.TestPost{
			bot.TestPost{
				UserID:  "test-user3",
				Message: "###### You have been scheduled for test-rotation#1.\n@test-user scheduled test-rotation#1.\n\nTODO runbook/info URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user5",
				Message: "###### You have been scheduled for test-rotation#1.\n@test-user scheduled test-rotation#1.\n\nTODO runbook/info URL/channel",
			},
		}, poster.DirectPosts)
		poster.Reset()

		task := mustRunTask(t, SL, `/lotto task show test-rotation#1`)
		require.Equal(t, sl.TaskStateScheduled, task.State)

		mustRunMulti(t, SL, `
			/lotto task start test-rotation#1 
			/lotto task finish test-rotation#1 
			`)
		require.Equal(t, []bot.TestPost{
			bot.TestPost{
				UserID:  "test-user3",
				Message: "###### Your test-rotation#1 started!\n@test-user started test-rotation#1.\n\nTODO runbook URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user5",
				Message: "###### Your test-rotation#1 started!\n@test-user started test-rotation#1.\n\nTODO runbook URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user3",
				Message: "###### Your test-rotation#1 finished!\n@test-user finished test-rotation#1.\n\nTODO runbook URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user5",
				Message: "###### Your test-rotation#1 finished!\n@test-user finished test-rotation#1.\n\nTODO runbook URL/channel",
			},
		}, poster.DirectPosts)
		poster.Reset()
	})

	t.Run("schedule updates calendar", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		poster := &bot.TestPoster{}
		SL, _ := getTestSLWithPoster(t, ctrl, poster)

		checkCal := func(id string, n int, t1, t2, t3 string) {
			user := mustRunUser(t, SL, `/lotto user show `+id)
			require.Equal(t, 2, len(user.Calendar))
			require.Equal(t, []*sl.Unavailable{
				&sl.Unavailable{
					Interval: types.Interval{
						Start:  types.MustParseTime(t1),
						Finish: types.MustParseTime(t2),
					},
					Reason:     "task",
					TaskID:     types.ID(fmt.Sprintf("test-rotation#%v", n)),
					RotationID: "test-rotation",
				},
				&sl.Unavailable{
					Interval: types.Interval{
						Start:  types.MustParseTime(t2),
						Finish: types.MustParseTime(t3),
					},
					Reason:     "grace",
					TaskID:     types.ID(fmt.Sprintf("test-rotation#%v", n)),
					RotationID: "test-rotation",
				},
			}, user.Calendar)
		}

		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=shift --beginning=2020-03-01
			/lotto rotation set task test-rotation --grace=201h
			/lotto user join test-rotation @u
			/lotto task new shift test-rotation -n 0
			/lotto task fill test-rotation#0
			`)
		checkCal("@u", 0, "2020-03-01T08:00", "2020-03-08T08:00", "2020-03-16T17:00")

		// When task is scheduled doesn't impact the expected start
		mustRun(t, SL, `/lotto task schedule test-rotation#0 --now=2020-03-15`)
		checkCal("@u", 0, "2020-03-01T08:00", "2020-03-08T08:00", "2020-03-16T17:00")

		mustRun(t, SL, `/lotto task start test-rotation#0 --now=2020-04-01`)
		checkCal("@u", 0, "2020-04-01T07:00", "2020-04-08T07:00", "2020-04-16T16:00")
		mustRun(t, SL, `/lotto task finish test-rotation#0 --now=2020-04-19`)
		checkCal("@u", 0, "2020-04-01T07:00", "2020-04-19T07:00", "2020-04-27T16:00")
	})
}
