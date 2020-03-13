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
			/lotto rotation param shift test-rotation -s 2030-01-10 -p monthly`, &r)
		require.NoError(t, err)
		require.Equal(t, sl.ShiftMaker, r.TaskMaker.Type)
		require.Equal(t, "everyMonth", r.TaskMaker.ShiftPeriod.String())
		require.Equal(t, "2030-01-10T08:00:00Z", r.TaskMaker.ShiftStart.Format(time.RFC3339))
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
		require.Equal(t, sl.TicketMaker, r.TaskMaker.Type)
	})

	t.Run("max", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param max -k webapp-2 --count 2 test-rotation
			/lotto rotation param max -k webapp-3 --count 1 test-rotation
			/lotto rotation param max -k server-1 --count 3 test-rotation
			/lotto rotation param max -k webapp-3 --clear test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 1}, r.TaskMaker.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"server-◉": 3, "webapp-▣": 2}, r.TaskMaker.Limit.TestAsMap())
	})

	t.Run("min", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto rotation param min -k webapp-2 --count 2 test-rotation
			/lotto rotation param min -k webapp-3 --count 1 test-rotation
			/lotto rotation param min -k server --count 3 test-rotation
			/lotto rotation param min -k webapp-3 --clear test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 1, "server-◉": 3, "webapp-▣": 2}, r.TaskMaker.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{}, r.TaskMaker.Limit.TestAsMap())
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
		require.Equal(t, []types.ID{}, r.TaskMaker.Limit.IDs())
		require.Equal(t, []types.ID{sl.AnySkillLevel.AsID()}, r.TaskMaker.Require.IDs())
		require.Equal(t, int64(1), r.TaskMaker.Require.IntSet.Get(sl.AnySkillLevel.AsID()))

		err = runCommands(t, SL, `
			/lotto rotation param max -k * --count 3 test-rotation
			/lotto rotation param min -k web-1 --count 1 test-rotation
		`)
		require.NoError(t, err)

		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, map[types.ID]int64{"*-*": 3}, r.TaskMaker.Limit.TestAsMap())
		require.Equal(t, map[types.ID]int64{"*-*": 1, "web-◉": 1}, r.TaskMaker.Require.TestAsMap())
	})
}
