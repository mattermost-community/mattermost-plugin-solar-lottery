// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCard struct {
	ID   ID
	Data string
}

func (si testCard) GetID() ID { return si.ID }

type testIndexProto []testCard

func (p testIndexProto) Len() int                 { return len(p) }
func (p testIndexProto) GetAt(n int) IndexCard    { return p[n] }
func (p testIndexProto) SetAt(n int, v IndexCard) { p[n] = v.(testCard) }

func (p testIndexProto) InstanceOf() IndexCardArray {
	inst := make(testIndexProto, 0)
	return &inst
}
func (p *testIndexProto) Ref() interface{} { return &p }
func (p *testIndexProto) Resize(n int) {
	*p = make(testIndexProto, n)
}

func TestIndexJSON(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		in := NewIndex(IDIndexProto, ID("test1"), ID("test2"))

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, `["test1","test2"]`, string(data))

		out := NewIndex(IDIndexProto)
		err = json.Unmarshal(data, &out)
		require.NoError(t, err)

		var ain, aout idIndexProto
		in.TestAsArray(&ain)
		out.TestAsArray(&aout)
		require.EqualValues(t, ain, aout)
	})
	t.Run("structs", func(t *testing.T) {
		proto := &testIndexProto{}
		in := NewIndex(proto,
			testCard{
				ID:   "id2",
				Data: "data2",
			},
			testCard{
				ID:   "id3",
				Data: "data3",
			},
			testCard{
				ID:   "id1",
				Data: "data1",
			},
		)

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, `[{"ID":"id2","Data":"data2"},{"ID":"id3","Data":"data3"},{"ID":"id1","Data":"data1"}]`, string(data))

		out := NewIndex(proto)
		err = json.Unmarshal(data, &out)
		require.NoError(t, err)

		var ain, aout testIndexProto
		in.TestAsArray(&ain)
		out.TestAsArray(&aout)
		require.EqualValues(t, ain, aout)
	})
}
