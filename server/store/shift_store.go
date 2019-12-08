// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type ShiftStore interface {
	LoadShift(rotationName string, number int) (*Shift, error)
	StoreShift(shift *Shift) error
}

type Shift struct {
	PluginVersion string

	// Mandatory attributes
	Number            int
	Start             string
	End               string
	RotationName      string
	MattermostUserIDs UserIDList `json:",omitempty"`
}

func NewShift(rotationName string, number int) *Shift {
	return &Shift{
		RotationName:      rotationName,
		Number:            number,
		MattermostUserIDs: UserIDList{},
	}
}

func (s *pluginStore) LoadShift(rotationName string, number int) (*Shift, error) {
	key := fmt.Sprintf("%v-%v", rotationName, number)
	shift := Shift{}
	err := kvstore.LoadJSON(s.shiftKV, key, &shift)
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (s *pluginStore) StoreShift(shift *Shift) error {
	key := fmt.Sprintf("%v-%v", shift.RotationName, shift.Number)
	err := kvstore.StoreJSON(s.shiftKV, key, shift)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Shift": shift,
	}).Debugf("Stored shift")
	return nil
}
