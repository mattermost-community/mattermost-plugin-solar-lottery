// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type filterf func(*solarLottery) error

func (sl *solarLottery) Filter(filters ...filterf) error {
	for _, filter := range filters {
		err := filter(sl)
		if err != nil {
			return err
		}
	}
	return nil
}

func withRotation(rotationID string) func(sl *solarLottery) error {
	return func(sl *solarLottery) error {
		return nil
	}
}

func withRotationExpanded(rotation *Rotation) func(sl *solarLottery) error {
	return func(sl *solarLottery) error {
		return sl.ExpandRotation(rotation)
	}
}

func withRotationIsNotArchived(rotation *Rotation) func(sl *solarLottery) error {
	return func(sl *solarLottery) error {
		if rotation.IsArchived {
			return errors.Errorf("rotation %s is archived", rotation.Markdown())
		}
		return nil
	}
}

func withActingUser(sl *solarLottery) error {
	if sl.actingUser != nil {
		return nil
	}
	user, _, err := sl.loadOrMakeStoredUser(sl.actingMattermostUserID)
	if err != nil {
		return err
	}
	sl.actingUser = user
	return nil
}

func withActingUserExpanded(sl *solarLottery) error {
	if sl.actingUser != nil && sl.actingUser.MattermostUser != nil {
		return nil
	}
	err := withActingUser(sl)
	if err != nil {
		return err
	}
	return sl.ExpandUser(sl.actingUser)
}

func withKnownSkills(sl *solarLottery) error {
	if sl.knownSkills != nil {
		return nil
	}

	skills, err := sl.SkillsStore.LoadKnownSkills()
	if err == store.ErrNotFound {
		sl.knownSkills = store.IDMap{}
		return nil
	}
	if err != nil {
		return err
	}
	sl.knownSkills = skills
	return nil
}

func withValidSkillName(skillName string) func(sl *solarLottery) error {
	return func(sl *solarLottery) error {
		err := sl.Filter(withKnownSkills)
		if err != nil {
			return err
		}
		for s := range sl.knownSkills {
			if s == skillName {
				return nil
			}
		}
		return errors.Errorf("skill %s is not found", skillName)
	}
}

func withKnownRotations(sl *solarLottery) error {
	if sl.knownRotations != nil {
		return nil
	}

	rr, err := sl.RotationStore.LoadKnownRotations()
	if err != nil {
		if err == store.ErrNotFound {
			rr = store.IDMap{}
		} else {
			return err
		}
	}

	sl.knownRotations = rr
	return nil
}

func withMattermostUsersExpanded(mattermostUsernames string) func(sl *solarLottery) error {
	return func(sl *solarLottery) error {
		users, err := sl.LoadMattermostUsers(mattermostUsernames)
		if err != nil {
			return err
		}

		sl.users = users
		return nil
	}
}
