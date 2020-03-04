// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import "encoding/json"

type IntValue struct {
	ID    ID
	Value int64
}

func NewIntValue(id ID, value int64) IntValue {
	return IntValue{
		ID:    id,
		Value: value,
	}
}

func (ii IntValue) GetID() ID { return ii.ID }

type intArrayProto []IntValue

func (p intArrayProto) Len() int             { return len(p) }
func (p intArrayProto) GetAt(n int) Value    { return p[n] }
func (p intArrayProto) SetAt(n int, v Value) { p[n] = v.(IntValue) }
func (p *intArrayProto) Ref() interface{}    { return p }
func (p intArrayProto) InstanceOf() ValueArray {
	inst := make(intArrayProto, 0)
	return &inst
}
func (p *intArrayProto) Resize(n int) {
	*p = make(intArrayProto, n)
}

type IntSet struct {
	ValueSet
}

func NewIntSet(vv ...IntValue) *IntSet {
	i := &IntSet{
		ValueSet: *NewValueSet(&intArrayProto{}),
	}
	for _, v := range vv {
		i.Set(v.ID, v.Value)
	}
	return i
}

func (set *IntSet) Set(id ID, v int64) {
	set.ValueSet.Set(NewIntValue(id, v))
}

func (set *IntSet) Get(id ID) int64 {
	v := set.ValueSet.Get(id)
	if v == nil {
		return 0
	}
	return v.(IntValue).Value
}

func (set *IntSet) MarshalJSON() ([]byte, error) {
	m := map[ID]int64{}
	for _, id := range set.ids {
		m[id] = set.Get(id)
	}
	return json.Marshal(m)
}

func (set *IntSet) UnmarshalJSON(data []byte) error {
	m := map[ID]int64{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	n := NewIntSet()
	*set = *n
	for k, v := range m {
		set.Set(k, v)
	}
	return nil
}
