// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

type stringIdentifiable string

func (si stringIdentifiable) GetID() string        { return string(si) }
func (si stringIdentifiable) Clone(bool) Cloneable { return si }

type stringSetProto []stringIdentifiable

func (p stringSetProto) Len() int                 { return len(p) }
func (p stringSetProto) GetAt(n int) IndexCard    { return p[n] }
func (p stringSetProto) SetAt(n int, v IndexCard) { p[n] = v.(stringIdentifiable) }

func (p stringSetProto) InstanceOf() IndexCardArray {
	inst := make(stringSetProto, 0)
	return &inst
}
func (p *stringSetProto) Ref() interface{} { return p }
func (p *stringSetProto) Resize(n int) {
	*p = make(stringSetProto, n)
}

var StringSetProto = &stringSetProto{}

type StringSet struct {
	*Index
}

func NewStringSet(vv ...string) *StringSet {
	s := &StringSet{
		Index: NewIndex(&stringSetProto{}),
	}
	for _, v := range vv {
		s.Set(v)
	}
	return s
}

func (s *StringSet) Set(v string) {
	s.Index.Set(stringIdentifiable(v))
}
