// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type RotationsStore interface {
	LoadRotations() (map[string]*Rotation, error)
	StoreRotations(rotations map[string]*Rotation) error
}

type Rotation struct {
	PluginVersion     string
	MattermostUserIDs map[string]string
	MaxSize           int
	MinBetweenServe   int
	Name              string
	Needs             []Need
	Period            string
	Start             string
}

type Need struct {
	Name  string
	Count int
	Skill string
	Level Level
}

func (s *pluginStore) LoadRotations() (map[string]*Rotation, error) {
	rotations := map[string]*Rotation{}
	err := kvstore.LoadJSON(s.rotationsKV, "rotations", &rotations)
	if err != nil {
		return nil, err
	}
	return rotations, nil
}

func (s *pluginStore) StoreRotations(rotations map[string]*Rotation) error {
	err := kvstore.StoreJSON(s.rotationsKV, "rotations", rotations)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Rotations": rotations,
	}).Debugf("Stored rotations")
	return nil
}

type Level int

const (
	None = Level(iota)
	Beginner
	Intermediate
	Advanced
	Expert
)

func (l Level) String() string {
	switch l {
	case Beginner:
		return "beginner"
	case Intermediate:
		return "intermediate"
	case Advanced:
		return "advanced"
	case Expert:
		return "expert"
	}
	return "none"
}
