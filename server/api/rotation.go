// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Rotation struct {
	*store.Rotation

	StartTime time.Time
	Users     UserMap
}

func (rotation *Rotation) init(api *api) error {
	start, err := time.Parse(DateFormat, rotation.Start)
	if err != nil {
		return err
	}
	rotation.StartTime = start

	if rotation.Users == nil {
		rotation.Users = UserMap{}
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

	users, err := api.LoadStoredUsers(rotation.MattermostUserIDs)
	if err != nil {
		return err
	}
	err = api.ExpandUserMap(users)
	if err != nil {
		return err
	}
	rotation.Users = users

	return nil
}

func withRotation(rotationID string) func(api *api) error {
	return func(api *api) error {
		return nil
	}
}

func withRotationExpanded(rotation *Rotation) func(api *api) error {
	return func(api *api) error {
		return api.ExpandRotation(rotation)
	}
}

func withRotationIsNotArchived(rotation *Rotation) func(api *api) error {
	return func(api *api) error {
		if rotation.IsArchived {
			return errors.Errorf("rotation %s is archived", MarkdownRotation(rotation))
		}
		return nil
	}
}

func (rotation *Rotation) ChangeNeed(skill string, level Level, newNeed store.Need) {
	for i, need := range rotation.Needs {
		if need.Skill == skill && need.Level == int(level) {
			rotation.Needs[i] = newNeed
			return
		}
	}
	rotation.Needs = append(rotation.Needs, newNeed)
}

func (rotation *Rotation) DeleteNeed(skill string, level Level) error {
	for i, need := range rotation.Needs {
		if need.Skill == skill && need.Level == int(level) {
			newNeeds := append([]store.Need{}, rotation.Needs[:i]...)
			if i+1 < len(rotation.Needs) {
				newNeeds = append(newNeeds, rotation.Needs[i+1:]...)
			}
			rotation.Needs = newNeeds
			return nil
		}
	}
	return errors.Errorf("%s is not found in rotation %s", MarkdownSkillLevel(skill, level), MarkdownRotation(rotation))
}

func (api *api) deleteUsersFromRotation(users UserMap, rotation *Rotation) error {
	for _, user := range users {
		_, ok := rotation.MattermostUserIDs[user.MattermostUserID]
		if !ok {
			return errors.Errorf("%s is not found in rotation %s", MarkdownUser(user), MarkdownRotation(rotation))
		}

		delete(user.NextRotationShift, rotation.RotationID)
		_, err := api.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
		delete(rotation.MattermostUserIDs, user.MattermostUserID)
		api.messageLeftRotation(user, rotation)
		api.Logger.Debugf("removed %s from %s.", MarkdownUser(user), MarkdownRotation(rotation))
		return nil
	}

	return api.RotationStore.StoreRotation(rotation.Rotation)
}
