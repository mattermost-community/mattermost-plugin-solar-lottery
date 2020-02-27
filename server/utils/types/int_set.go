// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import "strconv"

type intIdentifiable int64

func (ii intIdentifiable) GetID() string       { return strconv.Itoa(int(ii)) }
func (ii intIdentifiable) Clone() Identifiable { return ii }

type intArrayProto []intIdentifiable

func (p intArrayProto) Len() int                    { return len(p) }
func (p intArrayProto) GetAt(n int) Identifiable    { return p[n] }
func (p intArrayProto) SetAt(n int, v Identifiable) { p[n] = v.(intIdentifiable) }

func (p intArrayProto) InstanceOf() IndexPrototype {
	inst := make(intArrayProto, 0)
	return &inst
}
func (p *intArrayProto) Ref() interface{} { return p }
func (p *intArrayProto) Resize(n int) {
	*p = make(intArrayProto, n)
}

type IntSet struct {
	Index
}

func NewIntSet(vv ...int64) *IntSet {
	s := &IntSet{
		Index: NewIndex(&intArrayProto{}),
	}
	for _, v := range vv {
		s.Set(v)
	}
	return s
}

func (s *IntSet) Set(v int64) {
	s.Index.Set(intIdentifiable(v))
}

func (s *IntSet) MarshalJSON() ([]byte, error) {
	return s.Index.MarshalJSON()
}

func (s *IntSet) UnmarshalJSON(data []byte) error {
	return s.Index.UnmarshalJSON(data)
}
