// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"time"
)

type Interval struct {
	Start  Time
	Finish Time
}

func (i *Interval) IsEmpty() bool {
	return i != nil && i.Start.Before(i.Finish.Time)
}

type RelInterval struct {
	Start  time.Duration
	Finish time.Duration
}

func (i *RelInterval) IsEmpty() bool {
	return i != nil && i.Start < i.Finish
}
