// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

const NotEmpty = "y"

type IDMap map[string]string

func (m IDMap) Clone() IDMap {
	n := IDMap{}
	for k, v := range m {
		n[k] = v
	}
	return n
}

type IntMap map[string]int

func (m IntMap) Clone() IntMap {
	n := IntMap{}
	for k, v := range m {
		n[k] = v
	}
	return n
}
