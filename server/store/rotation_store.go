// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type RotationStore interface {
	LoadRotation(rotationName string) (*Rotation, error)
	StoreRotation(rotation *Rotation) error
	DeleteRotation(rotationName string) error
}

type Rotation struct {
	PluginVersion     string
	Name              string
	MattermostUserIDs map[string]string
	MinBetweenServe   int
	Period            string
	Start             string
	MaxSize           int
	Needs             []Need
}

type Need struct {
	Name  string
	Count int
	Skill string
	Level Level
}

func (s *pluginStore) LoadRotation(rotationName string) (*Rotation, error) {
	rotation := Rotation{}
	err := kvstore.LoadJSON(s.userKV, rotationName, &rotation)
	if err != nil {
		return nil, err
	}
	return &rotation, nil
}

func (s *pluginStore) StoreRotation(rotation *Rotation) error {
	err := kvstore.StoreJSON(s.userKV, rotation.Name, rotation)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Rotation": rotation,
	}).Debugf("Stored rotation")
	return nil
}

func (s *pluginStore) DeleteRotation(rotationName string) error {
	err := s.userKV.Delete(rotationName)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"RotationName": rotationName,
	}).Debugf("Deleted rotation")
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
