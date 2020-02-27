// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type stringIdentifiable string

func (si stringIdentifiable) GetID() string       { return string(si) }
func (si stringIdentifiable) Clone() Identifiable { return si }

type structIdentifiable struct {
	ID   string
	Data string
}

func (si structIdentifiable) GetID() string       { return si.ID }
func (si structIdentifiable) Clone() Identifiable { return si }

type stringArrayProto []stringIdentifiable

func (p stringArrayProto) Len() int                    { return len(p) }
func (p stringArrayProto) GetAt(n int) Identifiable    { return p[n] }
func (p stringArrayProto) SetAt(n int, v Identifiable) { p[n] = v.(stringIdentifiable) }

func (p stringArrayProto) InstanceOf() IndexPrototype {
	inst := make(stringArrayProto, 0)
	return &inst
}
func (p *stringArrayProto) Ref() interface{} { return p }
func (p *stringArrayProto) Resize(n int) {
	*p = make(stringArrayProto, n)
}

type structArrayProto []structIdentifiable

func (p structArrayProto) Len() int                    { return len(p) }
func (p structArrayProto) GetAt(n int) Identifiable    { return p[n] }
func (p structArrayProto) SetAt(n int, v Identifiable) { p[n] = v.(structIdentifiable) }

func (p structArrayProto) InstanceOf() IndexPrototype {
	inst := make(structArrayProto, 0)
	return &inst
}
func (p *structArrayProto) Ref() interface{} { return &p }
func (p *structArrayProto) Resize(n int) {
	*p = make(structArrayProto, n)
}

func TestIndexJSON(t *testing.T) {
	t.Run("Marshal strings", func(t *testing.T) {
		proto := &stringArrayProto{}
		in := NewIndex(proto, stringIdentifiable("test1"), stringIdentifiable("test2"))

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, `["test1","test2"]`, string(data))

		out := NewIndex(proto)
		err = json.Unmarshal(data, &out)
		require.NoError(t, err)

		var ain, aout stringArrayProto
		in.AsArray(&ain)
		out.AsArray(&aout)
		require.EqualValues(t, ain, aout)
	})
	t.Run("Marshal structs", func(t *testing.T) {
		proto := &structArrayProto{}
		in := NewIndex(proto,
			structIdentifiable{
				ID:   "id2",
				Data: "data2",
			},
			structIdentifiable{
				ID:   "id3",
				Data: "data3",
			},
			structIdentifiable{
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

		var ain, aout structArrayProto
		in.AsArray(&ain)
		out.AsArray(&aout)
		require.EqualValues(t, ain, aout)
	})
}
