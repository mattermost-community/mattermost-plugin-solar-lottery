// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

func (api *api) ExpandUserMap(users UserMap) error {
	for _, user := range users {
		err := api.ExpandUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *api) ExpandRotation(rotation *Rotation) error {
	if !rotation.StartTime.IsZero() && len(rotation.Users) == len(rotation.MattermostUserIDs) {
		return nil
	}

	err := rotation.init(api)
	if err != nil {
		return err
	}

	if len(rotation.Users) == 0 {
		users, err := api.LoadStoredUsers(rotation.MattermostUserIDs)
		if err != nil {
			return err
		}
		err = api.ExpandUserMap(users)
		if err != nil {
			return err
		}
		rotation.Users = users
	}

	return nil
}

func (api *api) ExpandUser(user *User) error {
	if user.MattermostUser != nil {
		return nil
	}
	mattermostUser, err := api.PluginAPI.GetMattermostUser(user.MattermostUserID)
	if err != nil {
		return err
	}
	user.MattermostUser = mattermostUser
	return nil
}
