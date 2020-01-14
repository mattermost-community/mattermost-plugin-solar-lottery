// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func TestShiftNumberForTime(t *testing.T) {
	tests := []struct {
		name    string
		r       store.Rotation
		date    string
		want    int
		wantErr bool
	}{
		{
			name: "happy0w",
			r: store.Rotation{
				Period: EveryWeek,
				Start:  "2019-12-21",
			},
			date: "2019-12-24",
			want: 0,
		}, {
			name: "happy1w",
			r: store.Rotation{
				Period: EveryWeek,
				Start:  "2019-12-21",
			},
			date: "2020-01-01",
			want: 1,
		}, {
			name: "happy1m",
			r: store.Rotation{
				Period: EveryMonth,
				Start:  "2019-12-21",
			},
			date: "2020-01-20",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rotationStartTime, err := time.Parse(DateFormat, tt.r.Start)
			require.Nil(t, err)
			rotation := &Rotation{
				Rotation:  &tt.r,
				StartTime: rotationStartTime,
			}
			startTime, err := time.Parse(DateFormat, tt.date)
			require.Nil(t, err)
			got, err := rotation.ShiftNumberForTime(startTime)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
