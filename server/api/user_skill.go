// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"
)

func (api *api) updateUserSkill(user *User, skillName string, level int) (*User, error) {
	if user.SkillLevels[skillName] == level {
		// nothing to do
		api.Logger.Debugf("nothing to do for user %s, already has skill %s (%v)", MarkdownUser(user), skillName, level)
		return user, nil
	}

	if level == 0 {
		_, ok := user.SkillLevels[skillName]
		if !ok {
			return nil, errors.Errorf("%s does not have skill %s", MarkdownUser(user), skillName)
		}
		delete(user.SkillLevels, skillName)
	} else {
		user.SkillLevels[skillName] = level
	}

	user, err := api.storeUserWelcomeNew(user)
	if err != nil {
		return nil, err
	}
	api.Logger.Debugf("%s (%v) skill added to user %s", skillName, level, MarkdownUser(user))
	return user, nil
}

func withValidSkillName(skillName string) func(api *api) error {
	return func(api *api) error {
		err := api.Filter(withKnownSkills)
		if err != nil {
			return err
		}
		for _, s := range api.knownSkills {
			if s == skillName {
				return nil
			}
		}
		return errors.Errorf("skill %s is not found", skillName)
	}
}
