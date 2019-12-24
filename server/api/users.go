// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"

type UserMap map[string]*User

func (m UserMap) Clone() UserMap {
	users := UserMap{}
	for id, user := range m {
		users[id] = user
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
		err := api.expandUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func withMattermostUsersExpanded(mattermostUsernames string) func(api *api) error {
	return func(api *api) error {
		err := api.Filter(withActingUserExpanded)
		if err != nil {
			return err
		}

		if mattermostUsernames == "" {
			api.users = UserMap{
				api.actingMattermostUserID: api.actingUser,
			}
			return nil
		}

		users, err := api.LoadMattermostUsers(mattermostUsernames)
		if err != nil {
			return err
		}

		api.users = users
		return nil
	}
}
