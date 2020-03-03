// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
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
			/lotto task param shift test-rotation -s 2030-01-10 -p monthly`, &r)
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
			/lotto task param ticket test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, sl.TicketMaker, r.TaskMaker.Type)
	})

	t.Run("max", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto task param max -k webapp-2 --count 2 test-rotation
			/lotto task param max -k webapp-3 --count 1 test-rotation
			/lotto task param max -k server-1 --count 3 test-rotation
			/lotto task param max -k webapp-3 --clear test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, []string{"server-◉", "webapp-▣"}, r.TaskMaker.Max.TestIDs())
		require.Equal(t, int64(2), r.TaskMaker.Max.Get("webapp-▣"))
		require.Equal(t, int64(3), r.TaskMaker.Max.Get("server-◉"))
	})

	t.Run("min", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto task param min -k webapp-2 --count 2 test-rotation
			/lotto task param min -k webapp-3 --count 1 test-rotation
			/lotto task param min -k server-1 --count 3 test-rotation
			/lotto task param min -k webapp-3 --clear test-rotation
			`)
		require.NoError(t, err)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation show test-rotation`, &r)
		require.NoError(t, err)
		require.Equal(t, []string{"server-◉", "webapp-▣"}, r.TaskMaker.Min.TestIDs())
		require.Equal(t, int64(2), r.TaskMaker.Min.Get("webapp-▣"))
		require.Equal(t, int64(3), r.TaskMaker.Min.Get("server-◉"))
	})
}
