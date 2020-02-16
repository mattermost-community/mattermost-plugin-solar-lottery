// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package timeutils

import (
	"time"
)

type Interval struct {
	Start  Time
	Finish Time
}

type RelInterval struct {
	Start  time.Duration
	Finish time.Duration
}
