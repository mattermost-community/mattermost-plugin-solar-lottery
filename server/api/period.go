// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Period struct {
	value string
}

var _ pflag.Value = (*Period)(nil)

const (
	EveryWeek     = "1w"
	EveryTwoWeeks = "2w"
	EveryMonth    = "1m"
)

func (p *Period) String() string {
	return p.value
}

func (p *Period) Type() string {
	return "rotation_period"
}

func (p *Period) Set(in string) error {
	switch strings.ToLower(in) {
	case EveryWeek, "w", "week":
		p.value = EveryWeek
	case EveryTwoWeeks, "2weeks", "biweekly", "bi-weekly":
		p.value = EveryTwoWeeks
	case EveryMonth, "m", "month":
		p.value = EveryMonth
	default:
		return errors.Errorf("period must be `%s`, `%s` or `%s`", EveryWeek, EveryTwoWeeks, EveryMonth)
	}
	return nil
}
