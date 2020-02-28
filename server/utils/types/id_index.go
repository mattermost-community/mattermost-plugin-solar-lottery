// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
)

type ID string

func (id ID) GetID() ID            { return id }
func (id ID) Clone(bool) Cloneable { return id }

type idIndexProto []ID

func (p idIndexProto) Len() int                 { return len(p) }
func (p idIndexProto) GetAt(n int) IndexCard    { return p[n] }
func (p idIndexProto) SetAt(n int, v IndexCard) { p[n] = v.(ID) }

func (p idIndexProto) InstanceOf() IndexCardArray {
	inst := make(idIndexProto, 0)
	return &inst
}
func (p *idIndexProto) Ref() interface{} { return p }
func (p *idIndexProto) Resize(n int) {
	*p = make(idIndexProto, n)
}

var IDIndexProto = &idIndexProto{}

type IDIndex struct {
	*Index
}

func NewIDIndex(vv ...ID) *IDIndex {
	i := &IDIndex{
		Index: NewIndex(&idIndexProto{}),
	}
	for _, v := range vv {
		i.Set(v)
	}
	return i
}

func (i *IDIndex) Set(v ID) {
	i.Index.Set(v)
}

func (i *IDIndex) Clone(deep bool) Cloneable {
	c := *i
	c.Index = i.Index.Clone(deep).(*Index)
	return &c
}

func (i *IDIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.IDs())
}

func (i *IDIndex) UnmarshalJSON(data []byte) error {
	ids := []ID{}
	err := json.Unmarshal(data, &ids)
	if err != nil {
		return err
	}

	n := NewIDIndex(ids...)
	*i = *n
	return nil
}
