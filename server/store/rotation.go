// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"time"
)

type Rotation struct {
	PluginVersion string
	RotationID    string
	IsArchived    bool

	// Mandatory attributes
	Name string
	Type string

	// Optional attributes
	Pool IDMap `json:",omitempty"`

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
		Name: name,
		Pool: IDMap{},
	}
}

func (rotation *Rotation) Clone(deep bool) *Rotation {
	newRotation := *rotation
	if deep {
		newRotation.Pool = rotation.Pool.Clone()
	}
	return &newRotation
}
