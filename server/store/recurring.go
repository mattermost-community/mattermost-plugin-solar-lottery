// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"time"
)

type Recurring struct {
	PluginVersion string
	RecurringID   string
	Requires      Needs `json:",omitempty"`
	Limits        Needs `json:",omitempty"`

	CycleStart  utils.Time
	CyclePeriod string
	Rel         *utils.RelInterval `json:",omitempty"`
	Grace       time.Duration      `json:",omitempty"`

	Autopilot RecurringAutopilot `json:",omitempty"`
}

type RecurringAutopilot struct {
}
