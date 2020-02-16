// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package timeutils

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02T15:04"
)

type Time struct {
	time.Time
}

var _ json.Marshaler = (*Time)(nil)

func TimeNow() Time {
	return Time{
		Time: time.Now(),
	}
}

// String is in UTC, use LocalString for local time
func (t Time) String() string {
	s := t.UTC().Format(TimeFormat)
	return strings.TrimSuffix(s, "T00:00")
}

func (t Time) LocalString() string {
	s := t.Format(TimeFormat)
	return strings.TrimSuffix(s, "T00:00")
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	s := ""
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}
