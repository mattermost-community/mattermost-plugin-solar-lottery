// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type ShiftStore interface {
	LoadShift(rotationID string, shiftNumber int) (*Shift, error)
	StoreShift(rotationID string, shiftNumber int, shift *Shift) error
	DeleteShift(rotationID string, shiftNumber int) error
}

const (
	ShiftStatusOpen      = "open"
	ShiftStatusCommitted = "committed"
	ShiftStatusFinished  = "finished"
	ShiftStatusStarted   = "started"
)

type Shift struct {
	PluginVersion string

	// Mandatory attributes
	Status string
	Start  string
	End    string

	// Optional
	MattermostUserIDs IDMap `json:",omitempty"`
}

func NewShift(start, end string, mattermostUserIDs IDMap) *Shift {
	if mattermostUserIDs == nil {
		mattermostUserIDs = IDMap{}
	}
	return &Shift{
		Start:             start,
		End:               end,
		MattermostUserIDs: mattermostUserIDs,
	}
}

func (s *pluginStore) LoadShift(rotationID string, shiftNumber int) (*Shift, error) {
	key := fmt.Sprintf("%v-%v", rotationID, shiftNumber)
	shift := NewShift("", "", nil)
	err := kvstore.LoadJSON(s.shiftKV, key, shift)
	if err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *pluginStore) StoreShift(rotationID string, shiftNumber int, shift *Shift) error {
	key := fmt.Sprintf("%v-%v", rotationID, shiftNumber)
	err := kvstore.StoreJSON(s.shiftKV, key, shift)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Shift": shift,
	}).Debugf("store: Stored shift %s %v", rotationID, shiftNumber)
	return nil
}

func (s *pluginStore) DeleteShift(rotationID string, shiftNumber int) error {
	key := fmt.Sprintf("%v-%v", rotationID, shiftNumber)
	err := s.shiftKV.Delete(key)
	if err != nil {
		return err
	}
	s.Logger.Debugf("store: Deleted shift %s %v", rotationID, shiftNumber)
	return nil
}
