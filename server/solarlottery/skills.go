// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Skills interface {
	ListSkills() (store.IDMap, error)
	AddSkill(string) error
	DeleteSkill(string) error
}

func (sl *solarLottery) ListSkills() (store.IDMap, error) {
	err := sl.Filter(
		withKnownSkills,
		withActingUser,
	)
	if err != nil {
		return nil, err
	}
	return sl.knownSkills, nil
}

func (sl *solarLottery) AddSkill(skillName string) error {
	err := sl.Filter(
		withKnownSkills,
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.AddSkill",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"Skill":          skillName,
	})

	if sl.knownSkills[skillName] != "" {
		return ErrAlreadyExists

	}
	sl.knownSkills[skillName] = store.NotEmpty

	err = sl.SkillsStore.StoreKnownSkills(sl.knownSkills)
	if err != nil {
		return err
	}

	logger.Infof("%s added skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}

func (sl *solarLottery) DeleteSkill(skillName string) error {
	err := sl.Filter(
		withKnownSkills,
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.AddSkill",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"Skill":          skillName,
	})

	newSkills := store.IDMap{}
	for s := range sl.knownSkills {
		if s != skillName {
			newSkills[s] = store.NotEmpty
		}
	}
	if len(newSkills) == len(sl.knownSkills) {
		return errors.Errorf("skill %s is not found ", skillName)
	}

	err = sl.SkillsStore.StoreKnownSkills(newSkills)
	if err != nil {
		return err
	}
	logger.Infof("%s deleted skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}

func MarkdownSkillLevel(skillName string, level Level) string {
	return fmt.Sprintf("%s%s", Level(level).String(), skillName)
}
