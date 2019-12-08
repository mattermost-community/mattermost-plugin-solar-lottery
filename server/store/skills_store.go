// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type SkillsStore interface {
	LoadSkills() ([]string, error)
	StoreSkills([]string) error
}

func (s *pluginStore) LoadSkills() ([]string, error) {
	skills := []string{}
	err := kvstore.LoadJSON(s.skillsKV, "skills", &skills)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func (s *pluginStore) StoreSkills(skills []string) error {
	return kvstore.StoreJSON(s.skillsKV, "skills", skills)
}
