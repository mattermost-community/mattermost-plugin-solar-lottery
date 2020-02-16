// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

type RecurringStore interface {
	LoadRecurring(recurringID string) (*Recurring, error)
	StoreRecurring(shift *Recurring) error
	DeleteRecurring(recurringID string) error
}

func (s *pluginStore) LoadRecurring(recurringID string) (*Recurring, error) {
	recurring := Recurring{}
	err := kvstore.LoadJSON(s.recurringKV, recurringID, &recurring)
	if err != nil {
		return nil, err
	}
	return &recurring, nil
}

func (s *pluginStore) StoreRecurring(recurring *Recurring) error {
	err := kvstore.StoreJSON(s.recurringKV, recurring.RecurringID, recurring)
	if err != nil {
		return err
	}
	s.Logger.With(bot.LogContext{
		"Recurring": recurring,
	}).Debugf("store: Stored recurring %s", recurring.RecurringID)
	return nil
}

func (s *pluginStore) DeleteRecurring(recurringID string) error {
	err := s.recurringKV.Delete(recurringID)
	if err != nil {
		return err
	}
	s.Logger.Debugf("store: Deleted recurring %s %v", recurringID)
	return nil
}
