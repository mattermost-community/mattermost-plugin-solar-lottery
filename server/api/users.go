// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"

type UserMap map[string]*User

func (m UserMap) Clone(deep bool) UserMap {
	users := UserMap{}
	for id, user := range m {
		if deep {
			users[id] = user.Clone()
		} else {
			users[id] = user
		}
	}
	return users
}

func (m UserMap) IDMap() store.IDMap {
	ids := store.IDMap{}
	for id := range m {
		ids[id] = store.NotEmpty
	}
	return ids
}

func (api *api) ExpandUserMap(users UserMap) error {
	for _, user := range users {
		err := api.ExpandUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func withMattermostUsersExpanded(mattermostUsernames string) func(api *api) error {
	return func(api *api) error {
		users, err := api.LoadMattermostUsers(mattermostUsernames)
		if err != nil {
			return err
		}

		api.users = users
		return nil
	}
}
