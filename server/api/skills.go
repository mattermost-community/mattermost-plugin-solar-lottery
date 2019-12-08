// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"errors"
	"sort"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Skills interface {
	ListSkills() ([]string, error)
	AddSkill(string) ([]string, error)
	DeleteSkill(string) ([]string, error)
}

var ErrSkillAlreadyExists = errors.New("skill already exists")

func (api *api) ListSkills() ([]string, error) {
	return api.SkillsStore.LoadSkills()
}

func (api *api) AddSkill(add string) ([]string, error) {
	skills, err := api.SkillsStore.LoadSkills()
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}
	for _, s := range skills {
		if s == add {
			return nil, ErrSkillAlreadyExists
		}
	}
	skills = append(skills, add)
	sort.Strings(skills)

	err = api.SkillsStore.StoreSkills(skills)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func (api *api) DeleteSkill(skill string) ([]string, error) {
	skills, err := api.SkillsStore.LoadSkills()
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}
	newSkills := []string{}
	for _, s := range skills {
		if s != skill {
			newSkills = append(newSkills, s)
		}
	}
	if len(newSkills) == len(skill) {
		return nil, store.ErrNotFound
	}

	err = api.SkillsStore.StoreSkills(newSkills)
	if err != nil {
		return nil, err
	}
	return newSkills, nil
}
