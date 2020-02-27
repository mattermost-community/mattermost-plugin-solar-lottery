// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

type stringIdentifiable string

func (si stringIdentifiable) GetID() string        { return string(si) }
func (si stringIdentifiable) Clone(bool) Cloneable { return si }

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

type StringSet struct {
	Index
}

func NewStringSet(vv ...string) *StringSet {
	s := &StringSet{
		Index: NewIndex(&stringArrayProto{}),
	}
	for _, v := range vv {
		s.Set(v)
	}
	return s
}

func (s *StringSet) Set(v string) {
	s.Index.Set(stringIdentifiable(v))
}

func (s *StringSet) MarshalJSON() ([]byte, error) {
	return s.Index.MarshalJSON()
}

func (s *StringSet) UnmarshalJSON(data []byte) error {
	return s.Index.UnmarshalJSON(data)
}
