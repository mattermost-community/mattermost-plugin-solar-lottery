// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRotationList(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		sl, _ := getTestSL(t, ctrl)

		runCommands(t, sl, `
			/lotto rotation new test
			/lotto rotation new test-123
			/lotto rotation new test-345`)

		out := []string{}
		_, err := runJSONCommand(t, sl, `
			/lotto rotation list`, &out)
		require.NoError(t, err)
		require.Equal(t, []string{"test", "test-123", "test-345"}, out)
	})
}
