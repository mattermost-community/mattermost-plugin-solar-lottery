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
	time.Time // always in UTC
}

// var _ json.Marshaler = (*Time)(nil)
// var _ json.Unmarshaler = (*Time)(nil)
var _ pflag.Value = (*Time)(nil)

func NewTime(tt ...time.Time) Time {
	if len(tt) == 0 {
		return Time{
			Time: time.Now(),
		}
	}

	return Time{
		Time: tt[0],
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

	tt, err := time.ParseInLocation(TimeFormat, in, loc)
	if err != nil {
		tt, err = time.ParseInLocation(DateFormat, in, loc)
		if err != nil {
			return err
		}
	}
	t.Time = tt
	return nil
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
	t := NewTime()
	err := t.Set(in)
	if err != nil {
		panic(err.Error())
	}
	return t
}
