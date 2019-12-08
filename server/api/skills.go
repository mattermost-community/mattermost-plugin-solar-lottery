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
	ValidateSkill(string) error
	AddSkill(string) ([]string, error)
	DeleteSkill(string) ([]string, error)
}

var ErrSkillAlreadyExists = errors.New("skill already exists")

func (api *api) ListSkills() ([]string, error) {
	err := api.Filter(withSkills)
	if err != nil {
		return nil, err
	}
	return api.skills, nil
}

func (api *api) ValidateSkill(skill string) error {
	err := api.Filter(withSkills)
	if err != nil {
		return err
	}
	api.Errorf("<><> LOOKING FOR %s", skill)

	for _, s := range api.skills {
		if s == skill {
			return nil
		}
	}

	api.Errorf("<><> NOT FOUND")
	return store.ErrNotFound
}

func (api *api) AddSkill(add string) ([]string, error) {
	err := api.Filter(withSkills)
	if err != nil {
		return nil, err
	}

	for _, s := range api.skills {
		if s == add {
			return nil, ErrSkillAlreadyExists
		}
	}
	api.skills = append(api.skills, add)
	sort.Strings(api.skills)

	err = api.SkillsStore.StoreSkills(api.skills)
	if err != nil {
		return nil, err
	}
	return api.skills, nil
}

func (api *api) DeleteSkill(skill string) ([]string, error) {
	err := api.Filter(withSkills)
	if err != nil {
		return nil, err
	}

	newSkills := []string{}
	for _, s := range api.skills {
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

func withSkills(api *api) error {
	if api.skills != nil {
		return nil
	}

	skills, err := api.SkillsStore.LoadSkills()
	if err == store.ErrNotFound {
		api.skills = []string{}
		return nil
	}
	if err != nil {
		return err
	}
	api.skills = skills
	return nil
}
