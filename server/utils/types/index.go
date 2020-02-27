// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"sort"
)

type Cloneable interface {
	Clone(deep bool) Cloneable
}

type IndexCard interface {
	Cloneable
	GetID() string
}

type IndexSetter interface {
	Set(IndexCard)
}

type IndexGetter interface {
	Get(IndexCard)
}

type IndexCardArray interface {
	Len() int
	GetAt(int) IndexCard
	SetAt(int, IndexCard)
	InstanceOf() IndexCardArray
	Ref() interface{}
	Resize(int)
}

type Index struct {
	proto IndexCardArray
	keys  []string
	m     map[string]IndexCard
}

func NewIndex(proto IndexCardArray, vv ...IndexCard) *Index {
	i := &Index{
		keys:  []string{},
		m:     map[string]IndexCard{},
		proto: proto,
	}
	for _, v := range vv {
		i.Set(v)
	}
	return i
}

func (i *Index) Clone(deep bool) Cloneable {
	n := NewIndex(i.proto)
	for _, key := range i.keys {
		n.Set(i.m[key].Clone(deep).(IndexCard))
	}
	return n
}

func (i *Index) Contains(id string) bool {
	_, ok := i.m[id]
	return ok
}

func (i *Index) Delete(keyToDelete string) {
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

func (i *Index) Get(key string) IndexCard {
	return i.m[key]
}

func (i *Index) GetAt(n int) IndexCard {
	return i.m[i.keys[n]]
}

func (i *Index) Len() int {
	return len(i.keys)
}

func (i *Index) Keys() []string {
	n := make([]string, len(i.keys))
	copy(n, i.keys)
	return n
}

func (i *Index) Set(v IndexCard) {
	id := v.GetID()
	if !i.Contains(id) {
		i.keys = append(i.keys, id)
	}
	i.m[id] = v
}

func (i *Index) SetAt(n int, v IndexCard) {
	id := v.GetID()
	if !i.Contains(id) {
		i.keys = append(i.keys, id)
	}
	i.m[id] = v
}

func (i *Index) MarshalJSON() ([]byte, error) {
	proto := i.proto.InstanceOf()
	proto.Resize(len(i.keys))
	for n, key := range i.keys {
		proto.SetAt(n, i.m[key])
	}
	return json.Marshal(proto)
}

func (i *Index) UnmarshalJSON(data []byte) error {
	proto := i.proto.InstanceOf()
	err := json.Unmarshal(data, proto.Ref())
	if err != nil {
		return err
	}

	i.keys = []string{}
	i.m = map[string]IndexCard{}

	for n := 0; n < proto.Len(); n++ {
		i.Set(proto.GetAt(n))
	}
	return nil
}

func (i *Index) TestAsArray(out IndexCardArray) {
	out.Resize(len(i.keys))
	for n, key := range i.keys {
		out.SetAt(n, i.m[key])
	}
}

func (i *Index) TestSortedKeys() []string {
	n := i.Keys()
	sort.Strings(n)
	return n
}
