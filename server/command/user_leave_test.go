// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserLeave(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, `
			/lotto rotation new test-rotation
			/lotto user join test-rotation @id1 @id2 @id3 @id4
			`)

		users := mustRunUsersJoin(t, SL, `/lotto user leave test-rotation @id2 @id3 @id5`)
		require.Equal(t, []string{"id2", "id3"}, users.TestIDs())

		r := mustRunRotation(t, SL, `/lotto rotation show test-rotation`)
		require.Equal(t, []string{"id1", "id4"}, r.MattermostUserIDs.TestIDs())
	})
}
