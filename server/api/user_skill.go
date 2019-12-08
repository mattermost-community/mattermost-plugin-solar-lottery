// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (api *api) UpdateUserSkill(skill string, level int) (*store.User, error) {
	err := api.Filter(withSkills, withUser)
	if err != nil {
		return nil, err
	}
	err = api.ValidateSkill(skill)
	if err != nil {
		return nil, err
	}

	if api.user.SkillLevels[skill] == level {
		// nothing to do
		return api.user, nil
	}

	api.user.SkillLevels[skill] = level

	err = api.UserStore.StoreUser(api.user)
	if err != nil {
		return nil, err
	}
	return api.user, nil
}

func (api *api) DeleteUserSkill(skill string) (*store.User, error) {
	err := api.Filter(withSkills, withUser)
	if err != nil {
		return nil, err
	}
	_, ok := api.user.SkillLevels[skill]
	if !ok {
		return nil, store.ErrNotFound
	}

	delete(api.user.SkillLevels, skill)

	err = api.UserStore.StoreUser(api.user)
	if err != nil {
		return nil, err
	}
	return api.user, nil
}
