// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type RotationStore interface {
	LoadActiveRotations() (IDMap, error)
	StoreActiveRotations(IDMap) error
	LoadRotation(string) (*Rotation, error)
	DeleteRotation(string) error
	StoreRotation(*Rotation) error
}

func (s *pluginStore) LoadActiveRotations() (IDMap, error) {
	rotations := IDMap{}
	err := kvstore.LoadJSON(s.basicKV, ActiveRotationsKey, &rotations)
	if err != nil {
		return nil, err
	}
	return rotations, nil
}

func (s *pluginStore) LoadRotation(rotationID string) (*Rotation, error) {
	rotation := NewRotation("")
	err := kvstore.LoadJSON(s.rotationKV, rotationID, rotation)
	if err != nil {
		return nil, err
	}
	return rotation, nil
}

func (s *pluginStore) StoreActiveRotations(rotations IDMap) error {
	err := kvstore.StoreJSON(s.basicKV, ActiveRotationsKey, rotations)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Rotations": rotations,
	}).Debugf("store: Stored active rotations")
	return nil
}

func (s *pluginStore) StoreRotation(rotation *Rotation) error {
	err := kvstore.StoreJSON(s.rotationKV, rotation.RotationID, rotation)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Rotation": rotation,
	}).Debugf("store: Stored rotation %s", rotation.RotationID)
	return nil
}

func (s *pluginStore) DeleteRotation(rotationID string) error {
	err := s.rotationKV.Delete(rotationID)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{}).Debugf("store: Deleted rotation %s", rotationID)
	return nil
}
