// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Event struct {
	store.Event
	StartTime time.Time
	EndTime   time.Time
}

func parseEventDates(start, end string) (time.Time, time.Time, error) {
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
