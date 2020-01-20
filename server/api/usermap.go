// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

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

func (api *api) ExpandUserMap(users UserMap) error {
	for _, user := range users {
		err := api.ExpandUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m UserMap) MarkdownWithSkills() string {
	out := []string{}
	for _, user := range m {
		out = append(out, fmt.Sprintf("%s %s", user.Markdown(), user.MarkdownSkills()))
	}
	return strings.Join(out, ", ")
}

func (m UserMap) Markdown() string {
	out := []string{}
	for _, user := range m {
		out = append(out, user.Markdown())
	}
	return strings.Join(out, ", ")
}

func (m UserMap) IDMap() store.IDMap {
	ids := store.IDMap{}
	for id := range m {
		ids[id] = store.NotEmpty
	}
	return ids
}

func (api *api) addEventToUsers(users UserMap, event Event, persist bool) error {
	for _, user := range users {
		user.AddEvent(event)

		if persist {
			_, err := api.storeUserWelcomeNew(user)
			if err != nil {
				return errors.WithMessagef(err, "failed to update user %s", user.Markdown())
			}
		}
	}
	return nil
}
