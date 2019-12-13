// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type ShiftStore interface {
	LoadShift(rotationName string, shiftNumber int) (*Shift, error)
	StoreShift(shift *Shift) error
}

const (
	ShiftStatusScheduled  = "scheduled"
	ShiftStatusClosed     = "closed"
	ShiftStatusInProgress = "inprogress"
)

type Shift struct {
	PluginVersion string

	// Mandatory attributes
	ShiftNumber       int
	ShiftStatus       string
	Start             string
	End               string
	RotationName      string
	MattermostUserIDs UserIDList `json:",omitempty"`
}

func NewShift(rotationName string, shiftNumber int) *Shift {
	return &Shift{
		RotationName:      rotationName,
		ShiftNumber:       shiftNumber,
		MattermostUserIDs: UserIDList{},
	}
}

func (s *pluginStore) LoadShift(rotationName string, shiftNumber int) (*Shift, error) {
	key := fmt.Sprintf("%v-%v", rotationName, shiftNumber)
	shift := Shift{}
	err := kvstore.LoadJSON(s.shiftKV, key, &shift)
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (s *pluginStore) StoreShift(shift *Shift) error {
	key := fmt.Sprintf("%v-%v", shift.RotationName, shift.ShiftNumber)
	err := kvstore.StoreJSON(s.shiftKV, key, shift)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Shift": shift,
	}).Debugf("store: Stored shift")
	return nil
}
