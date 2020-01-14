// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (api *api) GetActingUser() (*User, error) {
	err := api.Filter(withActingUser)
	if err != nil {
		return nil, err
	}
	return api.actingUser, nil
}

func (api *api) LoadStoredUsers(ids store.IDMap) (UserMap, error) {
	users := UserMap{}
	for id := range ids {
		u, err := api.UserStore.LoadUser(id)
		if err != nil {
			return nil, err
		}
		users[u.MattermostUserID] = &User{
			User: u,
		}
	}
	return users, nil
}

func (api *api) LoadMattermostUsers(mattermostUsernames string) (UserMap, error) {
	err := api.Filter(withActingUserExpanded)
	if err != nil {
		return nil, err
	}

	if mattermostUsernames == "" {
		return UserMap{
			api.actingMattermostUserID: api.actingUser,
		}, nil
	}

	users := UserMap{}
	names := strings.Split(mattermostUsernames, ",")
	for _, name := range names {
		if strings.HasPrefix(name, "@") {
			name = name[1:]
		}
		mmuser, err := api.PluginAPI.GetMattermostUserByUsername(name)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load %s", name)
		}
		user, _, err := api.loadOrMakeStoredUser(mmuser.Id)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load User %s", name)
		}
		user.MattermostUser = mmuser
		users[mmuser.Id] = user
	}
	return users, nil
}
