// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type filterf func(*sl) error

func (sl *sl) Filter(filters ...filterf) error {
	for _, filter := range filters {
		err := filter(sl)
		if err != nil {
			return err
		}
	}
	return nil
}

func withRotation(rotationID string) func(sl *sl) error {
	return func(sl *sl) error {
		return nil
	}
}

func withRotationExpanded(rotation *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		return sl.ExpandRotation(rotation)
	}
}

func withRotationIsNotArchived(rotation *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		if rotation.IsArchived {
			return errors.Errorf("rotation %s is archived", rotation.Markdown())
		}
		return nil
	}
}

func withActingUser(sl *sl) error {
	if sl.actingUser != nil {
		return nil
	}
	user, _, err := sl.loadOrMakeUser(sl.actingMattermostUserID)
	if err != nil {
		return err
	}
	sl.actingUser = user
	return nil
}

func withActingUserExpanded(sl *sl) error {
	if sl.actingUser != nil && sl.actingUser.mattermostUser != nil {
		return nil
	}
	err := withActingUser(sl)
	if err != nil {
		return err
	}
	return sl.ExpandUser(sl.actingUser)
}

func withKnownSkills(sl *sl) error {
	if sl.knownSkills != nil {
		return nil
	}

	skills, err := sl.Store.Index(KeyKnownSkills).Load()
	if err == kvstore.ErrNotFound {
		sl.knownSkills = types.NewSet()
		return nil
	}
	if err != nil {
		return err
	}
	sl.knownSkills = skills
	return nil
}

func withValidSkillName(skillName string) func(sl *sl) error {
	return func(sl *sl) error {
		err := sl.Filter(withKnownSkills)
		if err != nil {
			return err
		}
		if !sl.knownSkills.Contains(skillName) {
			return errors.Errorf("skill %s is not found", skillName)
		}
		return nil
	}
}

func withActiveRotations(sl *sl) error {
	if sl.activeRotations != nil {
		return nil
	}

	rotations, err := sl.Store.Index(KeyActiveRotations).Load()
	if err == kvstore.ErrNotFound {
		sl.activeRotations = types.NewSet()
		return nil
	}
	if err != nil {
		return err
	}
	sl.activeRotations = rotations
	return nil
}
