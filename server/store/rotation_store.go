// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type RotationStore interface {
	LoadKnownRotations() (IDMap, error)
	StoreKnownRotations(IDMap) error
	LoadRotation(string) (*Rotation, error)
	DeleteRotation(rotationID string) error
	StoreRotation(*Rotation) error
}

type Rotation struct {
	PluginVersion string
	RotationID    string
	IsArchived    bool

	// Mandatory attributes
	Name   string
	Period string
	Start  string
	Type   string

	// Optional attributes
	Size              int   `json:",omitempty"`
	Grace             int   `json:",omitempty"`
	MattermostUserIDs IDMap `json:",omitempty"`
	Needs             Needs `json:",omitempty"`

	Autopilot RotationAutopilot `json:",omitempty"`
}

type RotationAutopilot struct {
	On          bool          `json:",omitempty"`
	StartFinish bool          `json:",omitempty"`
	Fill        bool          `json:",omitempty"`
	FillPrior   time.Duration `json:",omitempty"`
	Notify      bool          `json:",omitempty"`
	NotifyPrior time.Duration `json:",omitempty"`
}

func NewRotation(name string) *Rotation {
	return &Rotation{
		Name:              name,
		MattermostUserIDs: IDMap{},
		Needs:             Needs{},
	}
}

func (rotation *Rotation) Clone(deep bool) *Rotation {
	newRotation := *rotation
	if deep {
		newRotation.MattermostUserIDs = rotation.MattermostUserIDs.Clone()
		newRotation.Needs = append(Needs{}, rotation.Needs...)
	}
	return &newRotation
}

func (s *pluginStore) LoadKnownRotations() (IDMap, error) {
	rotations := IDMap{}
	err := kvstore.LoadJSON(s.basicKV, KnownRotationsKey, &rotations)
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

func (s *pluginStore) StoreKnownRotations(rotations IDMap) error {
	err := kvstore.StoreJSON(s.basicKV, KnownRotationsKey, rotations)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Rotations": rotations,
	}).Debugf("store: Stored known rotations")
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
