// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Period struct {
	value string
	dur   time.Duration
}

var _ pflag.Value = (*Period)(nil)

const (
	EveryWeek     = "everyWeek"
	EveryDay      = "everyDay"
	EveryTwoWeeks = "everyTwoWeeks"
	EveryMonth    = "everyMonth"
	EveryDuration = "everyDuration"
)

func (p *Period) String() string {
	return p.value
}

func (p *Period) Type() string {
	return "period"
}

func (p *Period) Set(in string) error {
	switch strings.ToLower(in) {
	case EveryDay, "d", "day", "daily":
		p.value = EveryDay
	case EveryWeek, "w", "week", "weekly":
		p.value = EveryWeek
	case EveryTwoWeeks, "2weeks", "biweekly", "bi-weekly":
		p.value = EveryTwoWeeks
	case EveryMonth, "m", "month":
		p.value = EveryMonth
	default:
		dur, err := time.ParseDuration(in)
		if err != nil {
			return errors.New(`period must be "daily", "workday", "weekly", "biweekly", "monthly", or a valid go duration`)
		}
		p.value = EveryDuration
		p.dur = dur
	}
	return nil
}

func (p *Period) ForTime(beginning, now Time) Interval {
	start := p.StartForTime(beginning, now)
	next := p.Next(start)
	return Interval{
		Start:  start,
		Finish: next,
	}
}

func (p *Period) StartForTime(start, now Time) Time {
	const (
		maxDay   = time.Hour * 24
		maxMonth = maxDay * 31
	)

	if now.Before(start.Time) {
		return start
	}

	n := 0
	delta := now.Sub(start.Time)
	days, months := 0, 0
	switch p.value {
	case EveryDuration:
		reduced := now.Add(-p.dur / 2)
		return Time{reduced.Round(p.dur)}

	case EveryDay:
		days = 1
		n = int(delta / maxDay)

	case EveryWeek:
		days = 7
		n = int(delta / (7 * maxDay))

	case EveryTwoWeeks:
		days = 14
		n = int(delta / (14 * maxDay))

	case EveryMonth:
		months = 1
		n = int(delta / maxMonth)

	default:
		return start
	}

	t := start.AddDate(0, months*n, days*n)
	if now.Before(t) {
		panic("<><>")
	}
	for {
		next := t.AddDate(0, months, days)
		if now.Before(next) {
			return Time{t}
		}
		t = next
	}
}

func (p *Period) Next(start Time) Time {
	days, months := 0, 0
	switch p.value {
	case EveryDuration:
		reduced := start.Add(-p.dur / 2)
		return Time{reduced.Round(p.dur)}

	case EveryDay:
		days = 1

	case EveryWeek:
		days = 7

	case EveryTwoWeeks:
		days = 14

	case EveryMonth:
		months = 1
	}

	return Time{start.AddDate(0, months, days)}
}
