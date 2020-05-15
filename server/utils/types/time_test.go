// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeSet(t *testing.T) {

	tests := []struct {
		in            string
		expectedError string
		// expectedString string
	}{
		{
			in: "2025-01-01T11:00",
			// expectedString: "2025-01-01T11:00",
		},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			tt := NewTime(time.Now()).In(EST)
			err := tt.Set(tc.in)
			if tc.expectedError != "" {
				require.Equal(t, tc.expectedError, fmt.Sprintf("%v", err))
				return
			}
			require.Equal(t, tc.in, tt.String())
			require.Equal(t, tc.in, tt.In(EST).String())
		})
	}
}
func TestTimeJSON(t *testing.T) {
	tests := []struct {
		time               time.Time
		expectedJSON       string
		expectedBackString string
		expectedLocal      string
	}{
		{
			time:               time.Date(2025, time.January, 1, 11, 0, 0, 0, PST),
			expectedJSON:       `2025-01-01T19:00:00Z`,
			expectedBackString: `2025-01-01T19:00`,
			expectedLocal:      `2025-01-01T14:00`,
		},
	}
	for _, tc := range tests {

		t.Run(tc.time.Format(time.RFC3339), func(t *testing.T) {
			tt := NewTime(tc.time).In(EST)

			data, err := json.Marshal(&tt)
			require.NoError(t, err)
			require.Equal(t, tc.expectedJSON, strings.Trim(string(data), `"`))

			back := NewTime(time.Now()).In(time.UTC)
			err = json.Unmarshal(data, &back)
			require.NoError(t, err)
			require.Equal(t, tc.expectedBackString, back.String())
			require.Equal(t, tc.expectedLocal, back.In(EST).String())
		})
	}
}
