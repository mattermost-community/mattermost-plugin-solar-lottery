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

func TestCommandTaskNewTicket(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, _ := getTestSL(t, ctrl)

		err := runCommands(t, SL, `
			/lotto rotation new test-rotation
			/lotto task param ticket test-rotation
			/lotto task param min -k webapp-2 --count 2 test-rotation
			/lotto task param max -k server-3 --count 1 test-rotation
			`)
		require.NoError(t, err)

		ts := time.Now()
		task := sl.NewTask("")
		_, err = runJSONCommand(t, SL, `
			/lotto task new ticket test-rotation --summary test-summary1
			`, &task)
		require.NoError(t, err)
		require.Equal(t, "test-plugin-version", task.PluginVersion)
		require.Equal(t, types.ID("test-rotation-1"), task.TaskID)
		require.Equal(t, types.ID("test-rotation"), task.RotationID)
		require.Equal(t, sl.TaskStatusPending, task.Status)
		require.True(t, task.Created.After(ts))
		require.Equal(t, "test-summary1", task.Summary)
		require.Equal(t, []string{"webapp-▣"}, task.Min.TestIDs())
		require.Equal(t, int64(2), task.Min.Get("webapp-▣"))
		require.Equal(t, []string{"server-◈"}, task.Max.TestIDs())
		require.Equal(t, int64(1), task.Max.Get("server-◈"))
	})
}
