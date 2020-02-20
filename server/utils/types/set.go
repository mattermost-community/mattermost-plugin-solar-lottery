// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import "encoding/json"

var notEmpty = struct{}{}

type Set struct {
	asArray []string
	asMap   map[string]struct{}
}

func NewSet() *Set {
	return &Set{
		asMap: map[string]struct{}{},
	}
}

func (s Set) Clone() *Set {
	n := NewSet()
	n.asMap = s.AsMap()
	n.asArray = s.AsArray()
	return n
}

func (s Set) Len() int {
	return len(s.asArray)
}

func (s Set) ForEach(f func(id string)) {
	for _, id := range s.asArray {
		f(id)
	}
}

func (s Set) ForEachWithError(f func(id string) error) error {
	for _, id := range s.asArray {
		if err := f(id); err != nil {
			return err
		}
	}
	return nil
}

func (s Set) ForEachWithBreak(f func(id string) bool) {
	for _, id := range s.asArray {
		if f(id) {
			return
		}
	}
}

func (s Set) AsMap() map[string]struct{} {
	m := map[string]struct{}{}
	for _, v := range s.asArray {
		m[v] = notEmpty
	}
	return m
}

func (s Set) AsArray() []string {
	n := make([]string, len(s.asArray))
	copy(n, s.asArray)
	return n
}

func (s Set) Contains(v string) bool {
	_, ok := s.asMap[v]
	return ok
}

func (s *Set) Add(v string) {
	if s.Contains(v) {
		return
	}
	s.asArray = append(s.asArray, v)
	s.asMap[v] = notEmpty
}

func (s *Set) Delete(v string) {
	if !s.Contains(v) {
		return
	}
	n := NewSet()
	for _, vv := range s.asArray {
		if vv != v {
			n.Add(v)
		}
	}
	s.asArray = n.asArray
	s.asMap = n.asMap
}

func (s *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.asArray)
}

func (s *Set) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.asArray)
	if err != nil {
		return err
	}
	for _, v := range s.asArray {
		s.asMap[v] = notEmpty
	}
	return nil
}
