// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func TestTaskTransition(t *testing.T) {
	t.Run("user messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		poster := &bot.TestPoster{}
		SL, _ := getTestSLWithPoster(t, ctrl, poster)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param ticket test-rotation
			/lotto user join test-rotation @test-user1 @test-user2 @test-user3 @test-user4 @test-user5 @test-user6
			`)
		require.NoError(t, err)
		poster.Reset()

		err = runCommands(t, SL, `
			/lotto task new ticket test-rotation --summary test-summary1
			/lotto task assign test-rotation#1 @test-user3 @test-user5
			/lotto task schedule test-rotation#1 
			`)
		require.NoError(t, err)
		require.Equal(t, []bot.TestPost{
			bot.TestPost{
				UserID:  "test-user3",
				Message: "###### You have been scheduled for test-rotation#1.\n@test-user-username scheduled test-rotation#1.\n\nTODO runbook/info URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user5",
				Message: "###### You have been scheduled for test-rotation#1.\n@test-user-username scheduled test-rotation#1.\n\nTODO runbook/info URL/channel",
			},
		}, poster.DirectPosts)
		poster.Reset()

		task := sl.NewTask("")
		runJSONCommand(t, SL, `
			/lotto task show test-rotation#1
		`, task)
		require.NoError(t, err)
		require.Equal(t, sl.TaskStateScheduled, task.State)

		err = runCommands(t, SL, `
			/lotto task start test-rotation#1 
			/lotto task finish test-rotation#1 
			`)
		require.NoError(t, err)
		require.Equal(t, []bot.TestPost{
			bot.TestPost{
				UserID:  "test-user3",
				Message: "###### Your test-rotation#1 started!\n@test-user-username started test-rotation#1.\n\nTODO runbook URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user5",
				Message: "###### Your test-rotation#1 started!\n@test-user-username started test-rotation#1.\n\nTODO runbook URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user3",
				Message: "###### Your test-rotation#1 finished!\n@test-user-username finished test-rotation#1.\n\nTODO runbook URL/channel",
			},
			bot.TestPost{
				UserID:  "test-user5",
				Message: "###### Your test-rotation#1 finished!\n@test-user-username finished test-rotation#1.\n\nTODO runbook URL/channel",
			},
		}, poster.DirectPosts)
		poster.Reset()
	})
}
