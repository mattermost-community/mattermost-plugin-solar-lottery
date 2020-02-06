// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Event struct {
	store.Event
	StartTime time.Time
	EndTime   time.Time
}

func NewShiftEvent(rotation *Rotation, shiftNumber int, shift *Shift) Event {
	s, _, _ := rotation.ShiftDatesForNumber(shiftNumber)
	_, e, _ := rotation.ShiftDatesForNumber(shiftNumber + rotation.Grace)

	return Event{
		Event: store.Event{
			Type:        store.EventTypeShift,
			Start:       s.Format(DateFormat),
			End:         e.Format(DateFormat),
			RotationID:  rotation.RotationID,
			ShiftNumber: shiftNumber,
		},
		StartTime: s,
		EndTime:   e,
	}
}

func NewPersonalEvent(startTime, endTime time.Time) Event {
	return Event{
		Event: store.Event{
			Type:  store.EventTypePersonal,
			Start: startTime.Format(DateFormat),
			End:   endTime.Format(DateFormat),
		},
		StartTime: startTime,
		EndTime:   endTime,
	}
}

func (event Event) Markdown() string {
	return fmt.Sprintf("%s: %s to %s",
		event.Type, event.Start, event.End)
}

func (sl *solarLottery) AddEvent(mattermostUsernames string, event Event) error {
	err := sl.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":            "sl.AddSkillToUsers",
		"ActingUsername":      sl.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"Event":               event,
	})

	err = sl.addEventToUsers(sl.users, event, true)
	if err != nil {
		return err
	}

	logger.Infof("%s added event %s to %s.",
		sl.actingUser.Markdown(), event.Markdown(), sl.users.MarkdownWithSkills())
	return nil
}

func (sl *solarLottery) DeleteEvents(mattermostUsernames string, startDate, endDate string) error {
	err := sl.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":            "sl.AddSkillToUsers",
		"ActingUsername":      sl.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"StartDate":           startDate,
		"EndDate":             endDate,
	})

	for _, user := range sl.users {
		intervalStart, intervalEnd, err := ParseDatePair(startDate, endDate)
		if err != nil {
			return err
		}

		_, err = user.OverlapEvents(intervalStart, intervalEnd, true)
		if err != nil {
			return errors.WithMessagef(err, "failed to remove events from %s to %s", startDate, endDate)
		}

		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return errors.WithMessagef(err, "failed to update user %s", user.Markdown())
		}
	}

	logger.Infof("%s deleted events from %s to %s from users %s.",
		sl.actingUser.Markdown(), startDate, endDate, sl.users.MarkdownWithSkills())
	return nil
}

func ParseDatePair(start, end string) (time.Time, time.Time, error) {
	s, err := time.Parse(DateFormat, start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	e, err := time.Parse(DateFormat, end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if s.After(e) {
		return time.Time{}, time.Time{}, errors.Errorf("event start %v after end %v", s, e)
	}
	return s, e, nil
}

type eventSorter struct {
	events []store.Event
	by     func(p1, p2 store.Event) bool
}

// Len is part of sort.Interface.
func (s *eventSorter) Len() int {
	return len(s.events)
}

// Swap is part of sort.Interface.
func (s *eventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *eventSorter) Less(i, j int) bool {
	return s.by(s.events[i], s.events[j])
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type eventsBy func(event1, event2 store.Event) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by eventsBy) Sort(events []store.Event) {
	ps := &eventSorter{
		events: events,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

func byStartDate(event1, event2 store.Event) bool {
	return event1.Start < event2.Start
}
