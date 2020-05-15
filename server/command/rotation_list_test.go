// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRotationList(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl, sl := defaultEnv(t)
		defer ctrl.Finish()

		mustRunMulti(t, sl, `
			/lotto rotation new test
			/lotto rotation new test-123
			/lotto rotation new test-345`)

		out := []string{}
		mustRunJSON(t, sl, `/lotto rotation list`, &out)
		require.Equal(t, []string{"test", "test-123", "test-345"}, out)
	})
}
