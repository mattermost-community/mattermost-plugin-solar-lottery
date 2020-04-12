// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestIceBreakerManual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	poster := &bot.TestPoster{}
	SL, _ := getTestSLWithPoster(t, ctrl, poster)

	err := runCommands(t, SL, `
		/lotto rotation new ICE
		/lotto rotation param shift ICE --period weekly --beginning 2020-03-05T09:00PST --duration 30m 
		/lotto rotation param min ICE --count 2
		/lotto user join ICE @test-user1 @test-user2 @test-user3 @test-user4 @test-user5
		/lotto user join ICE @test-user6 @test-user7 @test-user8 @test-user9 @test-user10
	`)
	require.NoError(t, err)
	poster.Reset()

	out, err := runTaskCreateCommand(t, SL, `/lotto task new shift ICE --number 1`)
	task := out.Task
	require.NoError(t, err)
	require.Equal(t, types.ID("ICE#1"), task.TaskID)
	require.Equal(t, sl.TaskStatePending, task.State)
	require.Equal(t, "2020-03-05T17:00", task.ExpectedStart.String())
	require.Equal(t, "UTC", task.ExpectedStart.Location().String())
	require.Equal(t, "30m0s", task.ExpectedDuration.String())
	require.Equal(t, map[types.ID]int64{"*-*": 2}, task.Require.TestAsMap())

	s, err := runCommand(t, SL, `/lotto task fill ICE#1`)
	require.NoError(t, err)
	require.Equal(t, "Auto-assigned @test-user7-username, @test-user10-username to ticket ICE#1.", s.String())

	err = runCommands(t, SL, `
		/lotto task schedule ICE#1
		/lotto task start ICE#1
		/lotto task new shift ICE -n 2
		/lotto task fill ICE#2
		/lotto task schedule ICE#2
		/lotto task finish ICE#1
	`)
	require.NoError(t, err)
	require.Equal(t, []bot.TestPost{
		bot.TestPost{UserID: "test-user7", Message: "###### You have been scheduled for ICE#1.\n@test-user-username scheduled ICE#1.\n\nTODO runbook/info URL/channel"},
		bot.TestPost{UserID: "test-user10", Message: "###### You have been scheduled for ICE#1.\n@test-user-username scheduled ICE#1.\n\nTODO runbook/info URL/channel"},
		bot.TestPost{UserID: "test-user7", Message: "###### Your ICE#1 started!\n@test-user-username started ICE#1.\n\nTODO runbook URL/channel"},
		bot.TestPost{UserID: "test-user10", Message: "###### Your ICE#1 started!\n@test-user-username started ICE#1.\n\nTODO runbook URL/channel"},
		bot.TestPost{UserID: "test-user6", Message: "###### You have been scheduled for ICE#2.\n@test-user-username scheduled ICE#2.\n\nTODO runbook/info URL/channel"},
		bot.TestPost{UserID: "test-user4", Message: "###### You have been scheduled for ICE#2.\n@test-user-username scheduled ICE#2.\n\nTODO runbook/info URL/channel"},
		bot.TestPost{UserID: "test-user7", Message: "###### Your ICE#1 finished!\n@test-user-username finished ICE#1.\n\nTODO runbook URL/channel"},
		bot.TestPost{UserID: "test-user10", Message: "###### Your ICE#1 finished!\n@test-user-username finished ICE#1.\n\nTODO runbook URL/channel"},
	}, poster.DirectPosts)
	poster.Reset()
}

func TestIceBreakerInsufficient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	SL, _ := getTestSLWithPoster(t, ctrl, nil)

	err := runCommands(t, SL, `
		/lotto rotation new ICE
		/lotto rotation param shift ICE --period weekly --beginning 2020-03-05T09:00PST --duration 30m 

		# 2 week grace period so we run out of users by ICE#3
		/lotto rotation param grace ICE --duration 336h
		
		/lotto rotation param min ICE --count 2
		/lotto user join ICE @test-user1 @test-user2 @test-user3 @test-user4 @test-user5
		/lotto task new shift ICE --number 1
		/lotto task fill ICE#1
		/lotto task schedule ICE#1
		/lotto task new shift ICE --number 2
		/lotto task fill ICE#2
		/lotto task schedule ICE#2
		/lotto task new shift ICE --number 3
		`)
	require.NoError(t, err)

	// This will fail: because of the 2 week grace
	err = runCommands(t, SL, `
		/lotto task fill ICE#3
	`)
	require.Error(t, err)
	require.Equal(t, "failed to fill task ICE#3: failed to fill ICE#3, filling need **1** *-*, unfilled needs **1** *-*: insufficient", err.Error())
}
