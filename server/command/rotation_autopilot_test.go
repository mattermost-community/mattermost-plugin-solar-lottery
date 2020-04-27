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

// - Set up a bi-weekly rotation
//   - a 3 week grace period
//   - 2 people/shift (need 6 people to function); give it 8 people
//   - Beginning of rotation: 2020-01-05T09:30
//   - Autopilot params: create-prior=800h schedule-prior=100h remind-start-prior=24h remind-finish-prior=24h
// - Run Autopilot for each day between 1/1/2020 and 2/1/2020
// 		2020-01-01 - create TEST#0, TEST#1, TEST#2
// 		2020-01-02 - fill and schedule TEST#0
// 		2020-01-04 - start reminder: 2 users of TEST#0
// 		2020-01-05 - start TEST#0
// 		2020-01-15 - create TEST#3
//		2020-01-16 - fill and schedule TEST#1
//		2020-01-19 - finish reminder: 2 users of TEST#0; start reminder: 2 users of TEST#1
// 		2020-01-20 - finish TEST#0, start TEST#1
// 		2020-01-29 - create TEST#4
// 		2020-01-30 - fill and schedule: TEST#2
func TestRotationAutopilot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	SL, store := getTestSL(t, ctrl)

	mustRunMulti(t, SL, `
			# 3 week grace period, 2 people/shift need 6 people to function; give 8 people
			
			/lotto rotation new TEST --task-type=shift --beginning=2020-01-05T09:30 --period=biweekly --seed=873647632
			/lotto rotation set task TEST --grace 400h
			/lotto rotation set require TEST --count 2
			/lotto rotation set autopilot TEST --create --create-prior 800h --schedule --schedule-prior=100h --start-finish --remind-start --remind-start-prior=24h --remind-finish --remind-finish-prior=24h
			/lotto user join TEST @test-user1 @test-user2 @test-user3 @test-user4 @test-user5 @test-user6 @test-user7 @test-user8 --starting 2020-01-01
		`)

	check := func(now, expected string) {
		o, cmderr := run(t, SL, `/lotto rotation autopilot TEST --now=`+now)
		require.NoError(t, cmderr, now)
		require.Equal(t, expected, string(o), now)
	}

	checkNothing := func(now string) {
		check(now, `@test-user ran autopilot on TEST for `+now+`.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: nothing to do
  - start reminder: nothing to do
  - start: nothing to do`)
	}

	check(`2020-01-01T12:00`, `@test-user ran autopilot on TEST for 2020-01-01T12:00.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: created 3 shifts:
    - created shift TEST#0
    - created shift TEST#1
    - created shift TEST#2
  - fill and schedule: nothing to do
  - start reminder: nothing to do
  - start: nothing to do`)

	task := sl.Task{}
	err := store.Entity(sl.KeyTask).Load("TEST#0", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStatePending, task.State)
	require.Equal(t, types.Time{}, task.ActualStart)
	require.Equal(t, types.Time{}, task.ActualFinish)
	require.Equal(t, false, task.AutopilotRemindedStart)
	require.Equal(t, false, task.AutopilotRemindedFinish)
	require.Equal(t, 336*time.Hour, task.ExpectedDuration)
	require.Equal(t, types.MustParseTime("2020-01-05T17:30"), task.ExpectedStart)
	require.Equal(t, 400*time.Hour, task.Grace)

	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#1", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStatePending, task.State)
	require.Equal(t, types.Time{}, task.ActualStart)
	require.Equal(t, types.Time{}, task.ActualFinish)
	require.Equal(t, false, task.AutopilotRemindedStart)
	require.Equal(t, false, task.AutopilotRemindedFinish)
	require.Equal(t, 336*time.Hour, task.ExpectedDuration)
	require.Equal(t, types.MustParseTime("2020-01-19T17:30"), task.ExpectedStart)
	require.Equal(t, 400*time.Hour, task.Grace)

	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#2", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStatePending, task.State)
	require.Equal(t, types.Time{}, task.ActualStart)
	require.Equal(t, types.Time{}, task.ActualFinish)
	require.Equal(t, false, task.AutopilotRemindedStart)
	require.Equal(t, false, task.AutopilotRemindedFinish)
	require.Equal(t, 336*time.Hour, task.ExpectedDuration)
	require.Equal(t, types.MustParseTime("2020-02-02T17:30"), task.ExpectedStart)
	require.Equal(t, 400*time.Hour, task.Grace)

	err = store.Entity(sl.KeyTask).Load("TEST#3", &task)
	require.Error(t, err)

	check(`2020-01-02T12:00`, `@test-user ran autopilot on TEST for 2020-01-02T12:00.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: processed 1 tasks:
    - Auto-assigned @test-user6 (none), @test-user2 (none) to ticket TEST#0, transitioned TEST#0 to scheduled
  - start reminder: nothing to do
  - start: nothing to do`)

	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#0", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStateScheduled, task.State)
	require.Equal(t, types.Time{}, task.ActualStart)
	require.Equal(t, types.Time{}, task.ActualFinish)
	require.Equal(t, false, task.AutopilotRemindedStart)
	require.Equal(t, false, task.AutopilotRemindedFinish)
	require.Equal(t, 336*time.Hour, task.ExpectedDuration)
	require.Equal(t, types.MustParseTime("2020-01-05T17:30"), task.ExpectedStart)
	require.Equal(t, 400*time.Hour, task.Grace)
	require.Equal(t, []string{"test-user2", "test-user6"}, task.MattermostUserIDs.TestIDs())

	checkNothing(`2020-01-03T12:00`)

	check(`2020-01-04T12:00`, `@test-user ran autopilot on TEST for 2020-01-04T12:00.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: nothing to do
  - start reminder: messaged 2 users of 1 tasks
  - start: nothing to do`)

	check(`2020-01-05T12:00`, `@test-user ran autopilot on TEST for 2020-01-05T12:00.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: nothing to do
  - start reminder: nothing to do
  - started: 1 tasks`)
	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#0", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStateStarted, task.State)
	// Note that ActualStart matches --run --now time, not ExpectedStart
	require.Equal(t, types.MustParseTime("2020-01-05T20:00"), task.ActualStart)
	require.Equal(t, types.Time{}, task.ActualFinish)
	require.Equal(t, true, task.AutopilotRemindedStart)
	require.Equal(t, false, task.AutopilotRemindedFinish)

	checkNothing(`2020-01-06`)
	checkNothing(`2020-01-07`)
	checkNothing(`2020-01-08`)
	checkNothing(`2020-01-09`)
	checkNothing(`2020-01-10`)
	checkNothing(`2020-01-11`)
	checkNothing(`2020-01-12`)
	checkNothing(`2020-01-13`)
	checkNothing(`2020-01-14`)

	check(`2020-01-15`, `@test-user ran autopilot on TEST for 2020-01-15.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: created 1 shifts:
    - created shift TEST#3
  - fill and schedule: nothing to do
  - start reminder: nothing to do
  - start: nothing to do`)
	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#3", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStatePending, task.State)
	require.Equal(t, types.MustParseTime("2020-02-16T17:30"), task.ExpectedStart)
	err = store.Entity(sl.KeyTask).Load("TEST#4", &task)
	require.Error(t, err)

	check(`2020-01-16`,
		`@test-user ran autopilot on TEST for 2020-01-16.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: processed 1 tasks:
    - Auto-assigned @test-user7 (none), @test-user1 (none) to ticket TEST#1, transitioned TEST#1 to scheduled
  - start reminder: nothing to do
  - start: nothing to do`)

	checkNothing(`2020-01-17`)
	checkNothing(`2020-01-18`)

	check(`2020-01-19`, `@test-user ran autopilot on TEST for 2020-01-19.
  - finish reminder: messaged 2 users of 1 tasks
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: nothing to do
  - start reminder: messaged 2 users of 1 tasks
  - start: nothing to do`)

	check(`2020-01-20`, `@test-user ran autopilot on TEST for 2020-01-20.
  - finish reminder: nothing to do
  - finished: 1 tasks
  - create shift: nothing to do
  - fill and schedule: nothing to do
  - start reminder: nothing to do
  - started: 1 tasks`)
	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#0", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStateFinished, task.State)
	// Note that ActualStart matches --run --now time, not ExpectedStart
	require.Equal(t, types.MustParseTime("2020-01-05T20:00"), task.ActualStart)
	require.Equal(t, types.MustParseTime("2020-01-20T08:00"), task.ActualFinish)
	require.Equal(t, true, task.AutopilotRemindedStart)
	require.Equal(t, true, task.AutopilotRemindedFinish)
	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#1", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStateStarted, task.State)
	// Note that ActualStart matches --run --now time, not ExpectedStart
	require.Equal(t, types.MustParseTime("2020-01-20T08:00"), task.ActualStart)
	require.Equal(t, types.Time{}, task.ActualFinish)
	require.Equal(t, true, task.AutopilotRemindedStart)
	require.Equal(t, false, task.AutopilotRemindedFinish)

	checkNothing(`2020-01-20`)
	checkNothing(`2020-01-21`)
	checkNothing(`2020-01-22`)
	checkNothing(`2020-01-23`)
	checkNothing(`2020-01-24`)
	checkNothing(`2020-01-25`)
	checkNothing(`2020-01-26`)
	checkNothing(`2020-01-27`)
	checkNothing(`2020-01-28`)

	check(`2020-01-29`, `@test-user ran autopilot on TEST for 2020-01-29.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: created 1 shifts:
    - created shift TEST#4
  - fill and schedule: nothing to do
  - start reminder: nothing to do
  - start: nothing to do`)
	task = sl.Task{}
	err = store.Entity(sl.KeyTask).Load("TEST#4", &task)
	require.NoError(t, err)
	require.Equal(t, sl.TaskStatePending, task.State)
	require.Equal(t, types.MustParseTime("2020-03-01T17:30"), task.ExpectedStart)

	check(`2020-01-30`, `@test-user ran autopilot on TEST for 2020-01-30.
  - finish reminder: nothing to do
  - finish: nothing to do
  - create shift: nothing to do
  - fill and schedule: processed 1 tasks:
    - Auto-assigned @test-user5 (none), @test-user3 (none) to ticket TEST#2, transitioned TEST#2 to scheduled
  - start reminder: nothing to do
  - start: nothing to do`)

	checkNothing(`2020-01-31`)
	checkNothing(`2020-02-01`)
}
