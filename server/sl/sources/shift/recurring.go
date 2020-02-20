// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package shift

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const RecurringEnt = "recurring_"

type Recurring struct {
	PluginVersion string
	RecurringID   string
	Requires      []sl.Need `json:",omitempty"`
	Limits        []sl.Need `json:",omitempty"`

	CycleStart  types.Time
	CyclePeriod string
	Rel         *types.RelInterval `json:",omitempty"`
	Grace       time.Duration      `json:",omitempty"`
}
