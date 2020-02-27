// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

type idInt struct {
	ID string
	V  int64
}

func NewIDInt(id string, value int64) IndexCard {
	return idInt{
		ID: id,
		V:  value,
	}
}

func (ii idInt) GetID() string        { return ii.ID }
func (ii idInt) Clone(bool) Cloneable { return ii }

type intArrayProto []idInt

func (p intArrayProto) Len() int                 { return len(p) }
func (p intArrayProto) GetAt(n int) IndexCard    { return p[n] }
func (p intArrayProto) SetAt(n int, v IndexCard) { p[n] = v.(idInt) }

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

func NewIntSet() *IntIndex {
	return &IntIndex{
		Index: NewIndex(&intArrayProto{}),
	}
}

func (s *IntIndex) Set(id string, v int64) {
	s.Index.Set(NewIDInt(id, v))
}

func (s *IntIndex) Get(id string) int64 {
	return s.Index.Get(id).(idInt).V
}
