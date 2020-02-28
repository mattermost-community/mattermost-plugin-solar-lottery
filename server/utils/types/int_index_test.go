// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntIndex(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		in := NewIntIndex()
		in.Set("b", 2)
		in.Set("c", 1000)
		in.Set("a", 1)

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, `{"a":1,"b":2,"c":1000}`, string(data))

		out := NewIntIndex()
		err = json.Unmarshal(data, &out)
		require.NoError(t, err)

		outdata, err := json.Marshal(in)
		require.NoError(t, err)
		require.EqualValues(t, data, outdata)
	})
}
