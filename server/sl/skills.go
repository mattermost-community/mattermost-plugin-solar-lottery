// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Skills interface {
	ListKnownSkills() (*types.Set, error)
	AddKnownSkill(string) error
	DeleteKnownSkill(string) error
}

func (sl *sl) ListKnownSkills() (*types.Set, error) {
	err := sl.Filter(
		withKnownSkills,
		withActingUser,
	)
	if err != nil {
		return nil, err
	}
	return sl.knownSkills, nil
}

func (sl *sl) AddKnownSkill(skillName string) error {
	err := sl.Filter(
		withKnownSkills,
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

	err = sl.Store.Index(KeyKnownSkills).AddTo(skillName)
	if err != nil {
		return err
	}
	sl.knownSkills.Add(skillName)

	logger.Infof("%s added known skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}

func (sl *sl) DeleteKnownSkill(skillName string) error {
	err := sl.Filter(
		withKnownSkills,
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

	err = sl.Store.Index(KeyKnownSkills).DeleteFrom(skillName)
	if err != nil {
		return err
	}
	sl.knownSkills.Delete(skillName)

	logger.Infof("%s deleted skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}
