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

func (p *Period) ForNumber(beginning Time, num int) Time {
	days, months := 0, 0
	switch p.Period {
	case EveryDuration:
		reduced := beginning.Add(-p.Duration / 2)
		return NewTime(reduced.Round(p.Duration))

	case EveryDay:
		days = 1 * num

	case EveryWeek:
		days = 7 * num

	case EveryTwoWeeks:
		days = 14 * num

	case EveryMonth:
		months = 1 * num
	}

	return NewTime(beginning.AddDate(0, months, days))
}

func (p *Period) ForTime(beginning, forTime Time) (int, Time) {
	for n := -1; ; n++ {
		if forTime.Before(beginning.Time) {
			return n, beginning
		}
		beginning = p.ForNumber(beginning, 1)
	}
}
