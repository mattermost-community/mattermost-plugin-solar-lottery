// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestRotationArchive(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)

		runCommands(t, SL, `
			/lotto rotation new test
			/lotto rotation new test-123
			/lotto rotation new test-345
			`)

		activeRotations, err := store.IDIndex(sl.KeyActiveRotations).Load()
		require.NoError(t, err)
		require.Equal(t, types.NewIDSet("test", "test-123", "test-345"), activeRotations)

		r := sl.NewRotation()
		_, err = runJSONCommand(t, SL, `
			/lotto rotation archive test-123`, &r)
		require.NoError(t, err)
		require.Equal(t, types.ID("test-123"), r.RotationID)
		require.True(t, r.IsArchived)

		activeRotations, err = store.IDIndex(sl.KeyActiveRotations).Load()
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-345"}, activeRotations.TestIDs())

		err = store.Entity(sl.KeyRotation).Load("test-123", &r)
		require.NoError(t, err)
		require.Equal(t, types.ID("test-123"), r.RotationID)
		require.True(t, r.IsArchived)

		rr := []string{}
		_, err = runJSONCommand(t, SL, `
			/lotto rotation list`, &rr)
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-345"}, rr)

		_, err = runCommand(t, SL, `
			/lotto rotation show test-123`)
		require.Equal(t, kvstore.ErrNotFound, err)
	})

}

func TestRotationDelete(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)

		runCommands(t, SL, `
			/lotto rotation new test
			/lotto rotation new test-123
			/lotto rotation new test-345
			`)

		var rotationID types.ID
		_, err := runJSONCommand(t, SL, `
			/lotto rotation debug-delete test-123`, &rotationID)
		require.NoError(t, err)
		require.Equal(t, types.ID("test-123"), rotationID)

		activeRotations, err := store.IDIndex(sl.KeyActiveRotations).Load()
		require.NoError(t, err)
		require.Equal(t, types.NewIDSet("test", "test-345"), activeRotations)

		var r sl.Rotation
		err = store.Entity(sl.KeyRotation).Load("test-123", &r)
		require.Equal(t, kvstore.ErrNotFound, err)

		rr := []string{}
		_, err = runJSONCommand(t, SL, `
			/lotto rotation list`, &rr)
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-345"}, rr)

		_, err = runCommand(t, SL, `
			/lotto rotation show test-123`)
		require.Equal(t, kvstore.ErrNotFound, err)
	})

}
