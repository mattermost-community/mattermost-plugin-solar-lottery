// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
