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

func TestCommandTaskNewShift(t *testing.T) {
	t.Run("happy simple", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param shift test-rotation --period weekly --beginning 2020-03-03
			/lotto rotation param min -s webapp-2 --count 2 test-rotation
			/lotto rotation param max -s server-3 --count 1 test-rotation
			`)
		require.NoError(t, err)

		out, err := runTaskCreateCommand(t, SL, `/lotto task new shift test-rotation --number 1`)
		task := out.Task
		require.NoError(t, err)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "", task.Summary)
		require.Equal(t, "", task.Description)
		require.Equal(t, "2020-03-03T08:00", task.ExpectedStart.String())
		require.Equal(t, "UTC", task.ExpectedStart.Location().String())
		// The shift is created in UTC, so it's 168 hours, not shortened for
		// daylight savings time (the user is in PST).
		require.Equal(t, "168h0m0s", task.ExpectedDuration.String())
		require.Equal(t, map[types.ID]int64{"*-*": 1, "webapp-▣": 2}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◈": 1}, task.Limit.TestAsMap())

		r := &sl.Rotation{}
		_, err = runJSONCommand(t, SL, `/lotto rotation show test-rotation`, r)
		require.NoError(t, err)
		require.Equal(t, []string{"test-rotation#1"}, r.TaskIDs.TestIDs())

		out, err = runTaskCreateCommand(t, SL, `/lotto task new shift test-rotation -n 2`)
		require.NoError(t, err)
		require.Equal(t, types.ID("test-rotation#2"), out.Task.TaskID)
		require.Equal(t, "2020-03-10T08:00", out.Task.ExpectedStart.String())

		out, err = runTaskCreateCommand(t, SL, `/lotto task new shift test-rotation -n 3`)
		require.NoError(t, err)
		require.Equal(t, types.ID("test-rotation#3"), out.Task.TaskID)
		require.Equal(t, "2020-03-17T08:00", out.Task.ExpectedStart.String())

		s, err := runCommand(t, SL, `/lotto task new shift test-rotation -n 4`)
		require.NoError(t, err)
		require.Equal(t, "created shift test-rotation#4.", s.String())
	})

	t.Run("error already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param shift test-rotation --period weekly --beginning 2020-03-03
			/lotto task new shift test-rotation --number 2
			`)
		require.NoError(t, err)

		_, err = runCommand(t, SL, `/lotto task new shift test-rotation --number 2`)
		require.Error(t, err)
		require.Equal(t, "failed to make shift #2 (2020-03-10T08:00 to 2020-03-17T08:00): shift test-rotation#2 (2020-03-10T08:00 - 2020-03-17T08:00) already exists, state pending", err.Error())
	})
}
