// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"sort"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const (
	ReasonTask     = "task"
	ReasonGrace    = "grace"
	ReasonPersonal = "personal"
)

type Unavailable struct {
	types.Interval
	Reason string

	TaskID string
}

func NewUnavailable(reason string, interval types.Interval) *Unavailable {
	return &Unavailable{
		Reason:   reason,
		Interval: interval,
	}
}

type unavailableSorter struct {
	uu []*Unavailable
	by func(p1, p2 *Unavailable) bool
}

// Len is part of sort.Interface.
func (s *unavailableSorter) Len() int {
	return len(s.uu)
}

// Swap is part of sort.Interface.
func (s *unavailableSorter) Swap(i, j int) {
	s.uu[i], s.uu[j] = s.uu[j], s.uu[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *unavailableSorter) Less(i, j int) bool {
	return s.by(s.uu[i], s.uu[j])
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type unavailableBy func(u1, u2 *Unavailable) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by unavailableBy) Sort(uu []*Unavailable) {
	ps := &unavailableSorter{
		uu: uu,
		by: by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

func byStartDate(u1, u2 *Unavailable) bool {
	return u1.Start.Before(u2.Start.Time)
}
