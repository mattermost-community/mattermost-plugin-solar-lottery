// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (api *api) ChangeRotationNeed(r *store.Rotation, name, skill string, level, count int) {
	if r.Needs == nil {
		r.Needs = map[string]store.Need{}
	}
	need := r.Needs[name]
	need.Skill = skill
	need.Level = level
	need.Count = count
	r.Needs[name] = need
}

func (api *api) RemoveRotationNeed(r *store.Rotation, name string) error {
	_, ok := r.Needs[name]
	if !ok {
		return store.ErrNotFound
	}
	delete(r.Needs, name)
	return nil
}
