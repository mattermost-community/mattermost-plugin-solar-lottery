// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import "encoding/json"

type IDInt struct {
	ID ID
	V  int64
}

func NewIDInt(id ID, value int64) IDInt {
	return IDInt{
		ID: id,
		V:  value,
	}
}

func (ii IDInt) GetID() ID            { return ii.ID }
func (ii IDInt) Clone(bool) Cloneable { return ii }

type intArrayProto []IDInt

func (p intArrayProto) Len() int                 { return len(p) }
func (p intArrayProto) GetAt(n int) IndexCard    { return p[n] }
func (p intArrayProto) SetAt(n int, v IndexCard) { p[n] = v.(IDInt) }

func (p intArrayProto) InstanceOf() IndexCardArray {
	inst := make(intArrayProto, 0)
	return &inst
}
func (p *intArrayProto) Ref() interface{} { return p }
func (p *intArrayProto) Resize(n int) {
	*p = make(intArrayProto, n)
}

type IntIndex struct {
	*Index
}

func NewIntIndex(vv ...IDInt) *IntIndex {
	i := &IntIndex{
		Index: NewIndex(&intArrayProto{}),
	}
	for _, v := range vv {
		i.Set(v.ID, v.V)
	}
	return i
}

func (index *IntIndex) Set(id ID, v int64) {
	index.Index.Set(NewIDInt(id, v))
}

func (index *IntIndex) Get(id ID) int64 {
	v := index.Index.Get(id)
	if v == nil {
		return 0
	}
	return v.(IDInt).V
}

func (index *IntIndex) Clone(deep bool) Cloneable {
	c := *index
	c.Index = index.Index.Clone(deep).(*Index)
	return &c
}

func (index *IntIndex) MarshalJSON() ([]byte, error) {
	m := map[ID]int64{}
	for _, id := range index.ids {
		m[id] = index.Get(id)
	}
	return json.Marshal(m)
}

func (index *IntIndex) UnmarshalJSON(data []byte) error {
	m := map[ID]int64{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	n := NewIntIndex()
	*index = *n
	for k, v := range m {
		index.Set(k, v)
	}
	return nil
}
