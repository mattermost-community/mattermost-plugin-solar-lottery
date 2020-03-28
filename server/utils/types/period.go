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
	Period   string
	Duration time.Duration `json:",omitempty"`
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
	return p.Period
}

func (p *Period) Type() string {
	return "period"
}

func (p *Period) Set(in string) error {
	switch strings.ToLower(in) {
	case EveryDay, "d", "day", "daily":
		p.Period = EveryDay
	case EveryWeek, "w", "week", "weekly":
		p.Period = EveryWeek
	case EveryTwoWeeks, "2weeks", "biweekly", "bi-weekly":
		p.Period = EveryTwoWeeks
	case EveryMonth, "m", "month", "monthly":
		p.Period = EveryMonth
	default:
		Duration, err := time.ParseDuration(in)
		if err != nil {
			return errors.New(`period must be "daily", "workday", "weekly", "biweekly", "monthly", or a valid go duration`)
		}
		p.Period = EveryDuration
		p.Duration = Duration
	}
	return nil
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
	switch p.Period {
	case EveryDuration:
		reduced := now.Add(-p.Duration / 2)
		return NewTime(reduced.Round(p.Duration))

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
		panic("<><> TODO")
	}
	for {
		next := t.AddDate(0, months, days)
		if now.Before(next) {
			return NewTime(t)
		}
		t = next
	}
}

func (p *Period) Next(start Time) Time {
	days, months := 0, 0
	switch p.Period {
	case EveryDuration:
		reduced := start.Add(-p.Duration / 2)
		return NewTime(reduced.Round(p.Duration))

	case EveryDay:
		days = 1

	case EveryWeek:
		days = 7

	case EveryTwoWeeks:
		days = 14

	case EveryMonth:
		months = 1
	}

	return NewTime(start.AddDate(0, months, days))
}

func (p *Period) NumberForTime(start, now Time) int {
	if start.Before(now.Time) {
		return -1
	}

	n := 0
	for {
		if !start.Before(now.Time) {
			return n
		}
		start = p.Next(start)
		n++
	}
}
