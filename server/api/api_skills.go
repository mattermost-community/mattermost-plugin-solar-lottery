// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrSkillAlreadyExists = errors.New("skill already exists")

func (api *api) ListSkills() (store.IDMap, error) {
	err := api.Filter(
		withKnownSkills,
		withActingUser,
	)
	if err != nil {
		return nil, err
	}
	return api.knownSkills, nil
}

func (api *api) AddSkill(skillName string) error {
	err := api.Filter(
		withKnownSkills,
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.AddSkill",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"Skill":          skillName,
	})

	if api.knownSkills[skillName] != "" {
		return ErrSkillAlreadyExists

	}
	api.knownSkills[skillName] = store.NotEmpty

	err = api.SkillsStore.StoreKnownSkills(api.knownSkills)
	if err != nil {
		return err
	}

	logger.Infof("%s added skill %s.", api.MarkdownUser(api.actingUser), skillName)
	return nil
}

func (api *api) DeleteSkill(skillName string) error {
	err := api.Filter(
		withKnownSkills,
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.AddSkill",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"Skill":          skillName,
	})

	newSkills := store.IDMap{}
	for s := range api.knownSkills {
		if s != skillName {
			newSkills[s] = store.NotEmpty
		}
	}
	if len(newSkills) == len(api.knownSkills) {
		return errors.Errorf("skill %s is not found ", skillName)
	}

	err = api.SkillsStore.StoreKnownSkills(newSkills)
	if err != nil {
		return err
	}
	logger.Infof("%s deleted skill %s.", api.MarkdownUser(api.actingUser), skillName)
	return nil
}
