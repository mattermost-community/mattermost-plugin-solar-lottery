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
	DateFormat   = "2006-01-02"
	DateFormatTZ = "2006-01-02MST"
	TimeFormat   = "2006-01-02T15:04"
	TimeFormatTZ = "2006-01-02T15:04MST"
)

type Time struct {
	time.Time
}

var _ pflag.Value = (*Time)(nil)

func NewTime(t time.Time) Time {
	return Time{
		Time: t,
	}
}

func (t Time) In(l *time.Location) Time {
	return Time{
		Time: t.Time.In(l),
	}
}

func (t *Time) Type() string {
	return "time"
}

// Set implements pflag.Var, treats the parameter as UTC; to parse local
// times, use t.In(location) before calling fs.Parse().
func (t *Time) Set(in string) error {
	loc := t.Time.Location()
	if loc == nil {
		loc = time.UTC
	}

	var err error
	for _, format := range []string{TimeFormatTZ, TimeFormat, DateFormatTZ, DateFormat} {
		tt, err := time.ParseInLocation(format, in, loc)
		if err == nil {
			t.Time = tt
			return nil
		}
	}
	return err
}

// String is in UTC, use LocalString for local time
func (t Time) String() string {
	return strings.TrimSuffix(t.Format(TimeFormat), "T00:00")
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.UTC())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	s := ""
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	t.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	return nil
}

// Exposed for testing other packages
var EST, PST *time.Location

func init() {
	var err error
	EST, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(err.Error())
	}
	PST, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err.Error())
	}
}

func MustParseTime(in string) Time {
	t := Time{}
	err := t.Set(in)
	if err != nil {
		panic(err.Error())
	}
	return t
}
