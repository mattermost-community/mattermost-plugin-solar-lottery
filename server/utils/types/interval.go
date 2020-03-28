// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import "time"

type Interval struct {
	Start  Time
	Finish Time
}

func NewInterval(start, finish Time) Interval {
	return Interval{
		Start:  start,
		Finish: finish,
	}
}

func NewDurationInterval(start Time, d time.Duration) Interval {
	return Interval{
		Start:  start,
		Finish: NewTime(start.Add(d)),
	}
}

func MustParseInterval(start, finish string) Interval {
	return NewInterval(MustParseTime(start), MustParseTime(finish))
}

func (i *Interval) IsEmpty() bool {
	return i != nil && i.Start.Before(i.Finish.Time)
}

func (i *Interval) Overlaps(other Interval) bool {
	if other.Start.Time.Before(i.Start.Time) {
		other.Start = i.Start
	}
	if other.Finish.After(i.Finish.Time) {
		other.Finish = i.Finish
	}
	return !other.IsEmpty()
}
