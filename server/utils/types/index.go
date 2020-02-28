// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"sort"
)

type IndexCard interface {
	GetID() ID
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
	ids   []ID
	m     map[ID]IndexCard
}

func NewIndex(proto IndexCardArray, vv ...IndexCard) *Index {
	i := &Index{
		ids:   []ID{},
		m:     map[ID]IndexCard{},
		proto: proto,
	}
	for _, v := range vv {
		i.Set(v)
	}
	return i
}

func (index *Index) Contains(id ID) bool {
	_, ok := index.m[id]
	return ok
}

func (index *Index) Delete(toDelete ID) {
	if !index.Contains(toDelete) {
		return
	}

	for n, key := range index.ids {
		if key == toDelete {
			updated := index.ids[:n]
			if n+1 < len(index.ids) {
				updated = append(updated, index.ids[n+1:]...)
			}
			index.ids = updated
		}
	}
	delete(index.m, toDelete)
}

func (index *Index) Get(id ID) IndexCard {
	return index.m[id]
}

func (index *Index) GetAt(n int) IndexCard {
	return index.m[index.ids[n]]
}

func (index *Index) Len() int {
	return len(index.ids)
}

func (index *Index) IDs() []ID {
	n := make([]ID, len(index.ids))
	copy(n, index.ids)
	return n
}

func (index *Index) Set(v IndexCard) {
	id := v.GetID()
	if !index.Contains(id) {
		index.ids = append(index.ids, id)
	}
	index.m[id] = v
}

func (index *Index) SetAt(n int, v IndexCard) {
	id := v.GetID()
	if !index.Contains(id) {
		index.ids = append(index.ids, id)
	}
	index.m[id] = v
}

func (index *Index) MarshalJSON() ([]byte, error) {
	proto := index.proto.InstanceOf()
	proto.Resize(len(index.ids))
	for n, id := range index.ids {
		proto.SetAt(n, index.m[id])
	}
	return json.Marshal(proto)
}

func (index *Index) UnmarshalJSON(data []byte) error {
	proto := index.proto.InstanceOf()
	err := json.Unmarshal(data, proto.Ref())
	if err != nil {
		return err
	}

	index.ids = []ID{}
	index.m = map[ID]IndexCard{}

	for n := 0; n < proto.Len(); n++ {
		index.Set(proto.GetAt(n))
	}
	return nil
}

func (index *Index) TestAsArray(out IndexCardArray) {
	out.Resize(len(index.ids))
	for n, key := range index.ids {
		out.SetAt(n, index.m[key])
	}
}

func (index *Index) TestIDs() []string {
	n := []string{}
	for _, id := range index.IDs() {
		n = append(n, string(id))
	}
	sort.Strings(n)
	return n
}
