// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

type IntMap map[string]int64

func (m IntMap) Clone() IntMap {
	n := IntMap{}
	for k, v := range m {
		n[k] = v
	}
	return n
}
