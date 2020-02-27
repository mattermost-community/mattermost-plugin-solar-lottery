// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntSet(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		in := NewIntSet(1, 1000, 2)

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, `[1,1000,2]`, string(data))

		out := NewIntSet()
		err = json.Unmarshal(data, &out)
		require.NoError(t, err)
		require.EqualValues(t, in.SortedKeys(), out.SortedKeys())
	})
}
