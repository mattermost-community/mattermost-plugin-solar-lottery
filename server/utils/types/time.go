// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02T15:04"
)

type Time struct {
	time.Time
}

var _ json.Marshaler = (*Time)(nil)
var _ json.Unmarshaler = (*Time)(nil)
var _ pflag.Value = (*Time)(nil)

func Now() Time {
	// TODO  use MM time zone
	return Time{
		Time: time.Now().UTC(),
	}
}

func NewTime(t time.Time) Time {
	// TODO  use MM time zone
	return Time{
		Time: t.UTC(),
	}
}

func (t *Time) Type() string {
	return "time"
}

func (t *Time) Set(in string) error {
	// TODO  use MM time zone
	tt, err := time.Parse(time.RFC3339, in)
	if err != nil {
		return err
	}
	t.Time = tt.UTC()
	return nil
}

// String is in UTC, use LocalString for local time
func (t Time) String() string {
	s := t.UTC().Format(TimeFormat)
	return strings.TrimSuffix(s, "T00:00")
}

func (t Time) LocalString() string {
	s := t.Local().Format(TimeFormat)
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
	t.Time = parsedTime.UTC()
	return nil
}
