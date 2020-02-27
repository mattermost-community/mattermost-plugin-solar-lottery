// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSet(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		in := NewStringSet("test1", "test2")

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, `["test1","test2"]`, string(data))

		out := NewStringSet()
		err = json.Unmarshal(data, &out)
		require.NoError(t, err)
		require.EqualValues(t, in.SortedKeys(), out.SortedKeys())
	})
}
