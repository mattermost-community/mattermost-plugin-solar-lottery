// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/stretchr/testify/require"
)

func TestRotationSetAutopilot(t *testing.T) {
	ctrl, SL := defaultEnv(t)
	defer ctrl.Finish()

	mustRunMulti(t, SL, `/lotto rotation new TEST --beginning 2020-01-05T09:30 --seed=42`)

	r := mustRunRotation(t, SL, "/lotto rotation set autopilot TEST "+
		"--create --create-prior 500h "+
		"--schedule --schedule-prior=100h "+
		"--start-finish "+
		"--remind-start --remind-start-prior=24h "+
		"--remind-finish --remind-finish-prior=24h ")
	require.Equal(t, types.ID("TEST"), r.RotationID)
	require.Equal(t, sl.AutopilotSettings{
		Create:            true,
		CreatePrior:       500 * time.Hour,
		RemindFinish:      true,
		RemindFinishPrior: 24 * time.Hour,
		RemindStart:       true,
		RemindStartPrior:  24 * time.Hour,
		Schedule:          true,
		SchedulePrior:     100 * time.Hour,
		StartFinish:       true,
	}, r.AutopilotSettings)

	r = mustRunRotation(t, SL, "/lotto rotation set autopilot TEST --off")
	require.Equal(t, sl.AutopilotSettings{}, r.AutopilotSettings)
}

func TestRotationSetFill(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRun(t, SL,
			`/lotto rotation new test-rotation`)

		r := mustRunRotation(t, SL,
			`/lotto rotation set fill test-rotation --fuzz=2 --seed=1234`)
		require.Equal(t, int64(2), r.FillSettings.Fuzz)
		require.Equal(t, int64(1234), r.FillSettings.Seed)
	})
}

func TestTaskSet(t *testing.T) {
	t.Run("limit", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation set limit -s webapp-2 --count 2 test-rotation
			/lotto rotation set limit -s webapp-3 --count 1 test-rotation
			/lotto rotation set limit -s server-1 --count 3 test-rotation
			/lotto rotation set limit -s webapp-3 --clear test-rotation
			`)

		r := mustRunRotation(t, SL, `/lotto rotation show test-rotation`)
		require.Equal(t, map[types.ID]int64{"any": 1}, r.TaskSettings.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◉": 3, "webapp-▣": 2}, r.TaskSettings.Limit.TestAsMap())
	})

	t.Run("require", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation set require -s webapp-2 --count 2 test-rotation
			/lotto rotation set require -s webapp-3 --count 1 test-rotation
			/lotto rotation set require -s server --count 3 test-rotation
			/lotto rotation set require -s webapp-3 --clear test-rotation
			`)

		r := mustRunRotation(t, SL, `/lotto rotation show test-rotation`)
		require.Equal(t, map[types.ID]int64{"any": 1, "server-◉": 3, "webapp-▣": 2}, r.TaskSettings.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{}, r.TaskSettings.Limit.TestAsMap())
	})

	t.Run("require-limit-any", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			# new rotation defaults to ticket, min 1			
			/lotto rotation new test-rotation
			`)

		r := mustRunRotation(t, SL, `/lotto rotation show test-rotation`)
		require.Equal(t, []types.ID{}, r.TaskSettings.Limit.IDs())
		require.Equal(t, []types.ID{sl.AnySkillLevel.AsID()}, r.TaskSettings.Require.IDs())
		require.Equal(t, int64(1), r.TaskSettings.Require.IntSet.Get(sl.AnySkillLevel.AsID()))

		mustRunMulti(t, SL, `
			/lotto rotation set limit -s any --count 3 test-rotation
			/lotto rotation set require -s web-1 --count 1 test-rotation
		`)

		r = mustRunRotation(t, SL, `/lotto rotation show test-rotation`)
		require.Equal(t, map[types.ID]int64{"any": 3}, r.TaskSettings.Limit.TestAsMap())
		require.Equal(t, map[types.ID]int64{"any": 1, "web-◉": 1}, r.TaskSettings.Require.TestAsMap())
	})

	t.Run("duration and grace", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation 
			/lotto rotation set task test-rotation --grace=400h --duration=200h
			`)

		r := mustRunRotation(t, SL, `/lotto rotation show test-rotation`)
		require.Equal(t, 400*time.Hour, r.TaskSettings.Grace)
		require.Equal(t, 200*time.Hour, r.TaskSettings.Duration)
	})
}
