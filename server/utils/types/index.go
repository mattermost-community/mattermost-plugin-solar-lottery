// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"sort"
)

type Identifiable interface {
	GetID() string
	Clone() Identifiable
}

type IndexArray interface {
	Len() int
	GetAt(int) Identifiable
	SetAt(int, Identifiable)
}

type IndexSetter interface {
	Set(Identifiable)
}

type IndexPrototype interface {
	IndexArray
	InstanceOf() IndexPrototype
	Ref() interface{}
	Resize(int)
}

type Index interface {
	IndexArray
	json.Marshaler
	json.Unmarshaler

	// Clone() Index
	AsMap(IndexSetter)
	AsArray(IndexPrototype)
	Contains(id string) bool
	Delete(keyToDelete string)
	Get(string) Identifiable
	Keys() []string
	Set(Identifiable)
	SortedKeys() []string
}

type index struct {
	proto IndexPrototype
	keys  []string
	m     map[string]Identifiable
}

func NewIndex(proto IndexPrototype, vv ...Identifiable) Index {
	i := &index{
		keys:  []string{},
		m:     map[string]Identifiable{},
		proto: proto,
	}
	for _, v := range vv {
		i.Set(v)
	}
	return i
}

func (i *index) AsArray(out IndexPrototype) {
	out.Resize(len(i.keys))
	for n, key := range i.keys {
		out.SetAt(n, i.m[key])
	}
}

func (i *index) AsMap(out IndexSetter) {
	for _, v := range i.m {
		out.Set(v)
	}
}

// func (i *index) Clone() Index {
// 	n := NewIndex()
// 	n.m = s.Map()
// 	n.asArray = s.Array()
// 	return n
// }

func (i *index) Contains(id string) bool {
	_, ok := i.m[id]
	return ok
}

func (i *index) Delete(keyToDelete string) {
	if !i.Contains(keyToDelete) {
		return
	}

	for n, key := range i.keys {
		if key != keyToDelete {
			updated := i.keys[:n]
			if n+1 < len(i.keys) {
				updated = append(updated, i.keys[:n+1]...)
			}
			i.keys = updated
		}
	}
	delete(i.m, keyToDelete)
}

func (i *index) Get(key string) Identifiable {
	return i.m[key]
}

func (i *index) GetAt(n int) Identifiable {
	return i.m[i.keys[n]]
}

func (i *index) Len() int {
	return len(i.keys)
}

func (i *index) Keys() []string {
	n := make([]string, len(i.keys))
	copy(n, i.keys)
	return n
}

func (i *index) Set(v Identifiable) {
	id := v.GetID()
	if !i.Contains(id) {
		i.keys = append(i.keys, id)
	}
	i.m[id] = v
}

func (i *index) SetAt(n int, v Identifiable) {
	id := v.GetID()
	if !i.Contains(id) {
		i.keys = append(i.keys, id)
	}
	i.m[id] = v
}

func (i *index) SortedKeys() []string {
	n := i.Keys()
	sort.Strings(n)
	return n
}

func (i *index) MarshalJSON() ([]byte, error) {
	proto := i.proto.InstanceOf()
	proto.Resize(len(i.keys))
	for n, key := range i.keys {
		proto.SetAt(n, i.m[key])
	}
	return json.Marshal(proto)
}

func (i *index) UnmarshalJSON(data []byte) error {
	proto := i.proto.InstanceOf()
	err := json.Unmarshal(data, proto.Ref())
	if err != nil {
		return err
	}

	i.keys = []string{}
	i.m = map[string]Identifiable{}

	for n := 0; n < proto.Len(); n++ {
		i.Set(proto.GetAt(n))
	}
	return nil
}
