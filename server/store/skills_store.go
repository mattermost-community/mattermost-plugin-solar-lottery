// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type SkillsStore interface {
	LoadKnownSkills() (IDMap, error)
	StoreKnownSkills(IDMap) error
}

func (s *pluginStore) LoadKnownSkills() (IDMap, error) {
	skills := IDMap{}
	err := kvstore.LoadJSON(s.basicKV, KnownSkillsKey, &skills)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func (s *pluginStore) StoreKnownSkills(skills IDMap) error {
	err := kvstore.StoreJSON(s.basicKV, KnownSkillsKey, skills)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"skills": skills,
	}).Debugf("store: Stored known skills")
	return nil

}
