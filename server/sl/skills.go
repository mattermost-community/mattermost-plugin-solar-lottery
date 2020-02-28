// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Skills interface {
	ListKnownSkills() (*types.IDIndex, error)
	AddKnownSkill(types.ID) error
	DeleteKnownSkill(types.ID) error
}

func (sl *sl) ListKnownSkills() (*types.IDIndex, error) {
	err := sl.Filter(
		withKnownSkills,
		withActingUser,
	)
	if err != nil {
		return nil, err
	}
	return sl.knownSkills, nil
}

func (sl *sl) AddKnownSkill(skillName types.ID) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.AddKnownSkill",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"Skill":          skillName,
	})

	err = sl.Store.IDIndex(KeyKnownSkills).Set(skillName)
	if err != nil {
		return err
	}

	logger.Infof("%s added known skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}

func (sl *sl) DeleteKnownSkill(skillName types.ID) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.DeleteKnownSkill",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"Skill":          skillName,
	})

	err = sl.Store.IDIndex(KeyKnownSkills).Delete(skillName)
	if err != nil {
		return err
	}

	logger.Infof("%s deleted skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}
