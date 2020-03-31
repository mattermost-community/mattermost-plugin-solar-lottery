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

func TestCommandTaskParam(t *testing.T) {
	t.Run("shift", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation param shift test-rotation --beginning 2030-01-10 --period monthly`, &r)
		require.NoError(t, err)
		require.Equal(t, sl.TypeShift, r.Type)
		require.Equal(t, types.EveryMonth, r.ShiftPeriod.String())
		require.Equal(t, "2030-01-10T08:00:00Z", r.Beginning.Format(time.RFC3339))
	})

	t.Run("ticket", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
		`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation param ticket test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, sl.TypeTicket, r.Type)
	})

	t.Run("max", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param max -s webapp-2 --count 2 test-rotation
			/lotto rotation param max -s webapp-3 --count 1 test-rotation
			/lotto rotation param max -s server-1 --count 3 test-rotation
			/lotto rotation param max -s webapp-3 --clear test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 1}, r.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◉": 3, "webapp-▣": 2}, r.Limit.TestAsMap())
	})

	t.Run("min", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param min -s webapp-2 --count 2 test-rotation
			/lotto rotation param min -s webapp-3 --count 1 test-rotation
			/lotto rotation param min -s server --count 3 test-rotation
			/lotto rotation param min -s webapp-3 --clear test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 1, "server-◉": 3, "webapp-▣": 2}, r.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{}, r.Limit.TestAsMap())
	})

	t.Run("min-max-any", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			# new rotation defaults to ticket, min 1			
			/lotto rotation new test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, []types.ID{}, r.Limit.IDs())
		require.Equal(t, []types.ID{sl.AnySkillLevel.AsID()}, r.Require.IDs())
		require.Equal(t, int64(1), r.Require.IntSet.Get(sl.AnySkillLevel.AsID()))

		err = runCommands(t, SL, `
			/lotto rotation param max -s * --count 3 test-rotation
			/lotto rotation param min -s web-1 --count 1 test-rotation
		`)
		require.NoError(t, err)

		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 3}, r.Limit.TestAsMap())
		require.Equal(t, map[types.ID]int64{"*-*": 1, "web-◉": 1}, r.Require.TestAsMap())
	})
}
