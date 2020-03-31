// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntervalOverlaps(t *testing.T) {
	tests := []struct {
		name            string
		i1, i2          Interval
		expectedOverlap bool
	}{
		{
			name:            "superset",
			i1:              MustParseInterval("2025-01-01T11:00", "2025-01-02T11:00"),
			i2:              MustParseInterval("2024-01-01T11:00", "2026-01-02T11:00"),
			expectedOverlap: true,
		}, {
			name:            "overlap",
			i1:              MustParseInterval("2025-01-01T11:00", "2025-01-02T11:00"),
			i2:              MustParseInterval("2025-01-01T15:00", "2026-01-03T11:00"),
			expectedOverlap: true,
		}, {
			name:            "borderline",
			i1:              MustParseInterval("2024-01-01T00:00", "2025-01-01T11:00"),
			i2:              MustParseInterval("2024-01-01T11:00", "2026-01-02T11:00"),
			expectedOverlap: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedOverlap, tc.i1.Overlaps(tc.i2))
		})
	}
}
