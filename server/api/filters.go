// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type filterf func(*api) error

func (api *api) Filter(filters ...filterf) error {
	for _, filter := range filters {
		err := filter(api)
		if err != nil {
			return err
		}
	}
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
			return errors.Errorf("rotation %s is archived", rotation.Markdown())
		}
		return nil
	}
}

func withActingUser(api *api) error {
	if api.actingUser != nil {
		return nil
	}
	user, _, err := api.loadOrMakeStoredUser(api.actingMattermostUserID)
	if err != nil {
		return err
	}
	api.actingUser = user
	return nil
}

func withActingUserExpanded(api *api) error {
	if api.actingUser != nil && api.actingUser.MattermostUser != nil {
		return nil
	}
	err := withActingUser(api)
	if err != nil {
		return err
	}
	return api.ExpandUser(api.actingUser)
}

func withKnownSkills(api *api) error {
	if api.knownSkills != nil {
		return nil
	}

	skills, err := api.SkillsStore.LoadKnownSkills()
	if err == store.ErrNotFound {
		api.knownSkills = store.IDMap{}
		return nil
	}
	if err != nil {
		return err
	}
	api.knownSkills = skills
	return nil
}

func withValidSkillName(skillName string) func(api *api) error {
	return func(api *api) error {
		err := api.Filter(withKnownSkills)
		if err != nil {
			return err
		}
		for s := range api.knownSkills {
			if s == skillName {
				return nil
			}
		}
		return errors.Errorf("skill %s is not found", skillName)
	}
}

func withKnownRotations(api *api) error {
	if api.knownRotations != nil {
		return nil
	}

	rr, err := api.RotationStore.LoadKnownRotations()
	if err != nil {
		if err == store.ErrNotFound {
			rr = store.IDMap{}
		} else {
			return err
		}
	}

	api.knownRotations = rr
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
