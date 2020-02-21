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

func TestCommandRotationArchive(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)

		runCommands(t, SL, `
			/lotto rotation new test
			/lotto rotation new test-123
			/lotto rotation new test-345
			`)

		activeRotations, err := store.Index(sl.KeyActiveRotations).Load()
		require.NoError(t, err)
		require.Equal(t, types.NewSet("test", "test-123", "test-345"), activeRotations)

		r := &sl.Rotation{}
		_, err = runJSONCommand(t, SL, `
			/lotto rotation archive test-123`, &r)
		require.NoError(t, err)
		require.Equal(t, &sl.Rotation{
			RotationID: "test-123",
			IsArchived: true,
		}, r)

		activeRotations, err = store.Index(sl.KeyActiveRotations).Load()
		require.NoError(t, err)
		require.Equal(t, types.NewSet("test", "test-345"), activeRotations)

		err = store.Entity(sl.KeyRotation).Load("test-123", &r)
		require.NoError(t, err)
		require.Equal(t, &sl.Rotation{
			RotationID: "test-123",
			IsArchived: true,
		}, r)

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

func TestCommandRotationDelete(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)

		runCommands(t, SL, `
			/lotto rotation new test
			/lotto rotation new test-123
			/lotto rotation new test-345
			`)

		r := &sl.Rotation{}
		_, err := runJSONCommand(t, SL, `
			/lotto rotation debug-delete test-123`, &r)
		require.NoError(t, err)
		require.Equal(t, &sl.Rotation{
			RotationID: "test-123",
		}, r)

		activeRotations, err := store.Index(sl.KeyActiveRotations).Load()
		require.NoError(t, err)
		require.Equal(t, types.NewSet("test", "test-345"), activeRotations)

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
