// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsersJSON(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		users := NewUsers()

		data, err := json.Marshal(users)
		require.NoError(t, err)
		require.Equal(t, "[]", string(data))
	})
}
