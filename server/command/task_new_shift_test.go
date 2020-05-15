// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskNewShift(t *testing.T) {
	t.Run("happy simple", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new	test-rotation --task-type=shift --beginning 2020-03-03 --period weekly
			/lotto rotation set require	test-rotation -s webapp-2 --count 2 
			/lotto rotation set limit test-rotation -s server-3 --count 1 
			`)

		task := mustRunTaskCreate(t, SL,
			`/lotto task new shift test-rotation --number 0`)
		require.Equal(t, types.ID("test-rotation#0"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatePending, task.State)
		require.Equal(t, "", task.Summary)
		require.Equal(t, "", task.Description)
		require.Equal(t, "2020-03-03T08:00", task.ExpectedStart.String())
		require.Equal(t, "UTC", task.ExpectedStart.Location().String())
		// The shift is created in UTC, so it's 168 hours, not shortened for
		// daylight savings time (the user is in PST).
		require.Equal(t, "168h0m0s", task.ExpectedDuration.String())
		require.Equal(t, map[types.ID]int64{"any": 1, "webapp-▣": 2}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◈": 1}, task.Limit.TestAsMap())

		r := mustRunRotation(t, SL,
			`/lotto rotation show test-rotation`)
		require.Equal(t, []string{"test-rotation#0"}, r.TaskIDs.TestIDs())

		task = mustRunTaskCreate(t, SL, `/lotto task new shift test-rotation -n 1`)
		require.Equal(t, types.ID("test-rotation#1"), task.TaskID)
		require.Equal(t, "2020-03-10T08:00", task.ExpectedStart.String())

		task = mustRunTaskCreate(t, SL, `/lotto task new shift test-rotation -n 2`)
		require.Equal(t, types.ID("test-rotation#2"), task.TaskID)
		require.Equal(t, "2020-03-17T08:00", task.ExpectedStart.String())

		s := mustRun(t, SL, `/lotto task new shift test-rotation -n 3`)
		require.Equal(t, "created shift test-rotation#3", s.String())
	})

	t.Run("error already exists", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation --task-type=shift --beginning 2020-03-03 --period weekly
			/lotto rotation set fill test-rotation 
			/lotto task new shift test-rotation --number 2
			`)

		_, err := run(t, SL, `/lotto task new shift test-rotation --number 2`)
		require.Error(t, err)
		require.Equal(t, "failed to make shift #2 (2020-03-17T08:00 to 2020-03-24T08:00): shift test-rotation#2 (2020-03-17T08:00 - 2020-03-24T08:00) already exists, state pending", err.Error())
	})
}
