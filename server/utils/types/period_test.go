// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPeriod(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		p := &Period{
			Period: EveryWeek,
		}

		// Accounts for daylight savings time
		start := MustParseTime("2020-03-03").In(PST)
		next := p.Next(start, 1)
		require.Equal(t, "2020-03-09T16:00", next.String())
		require.Equal(t, "167h0m0s", next.Sub(start.Time).String())

		next2 := p.Next(next, 1)
		require.Equal(t, "2020-03-16T16:00", next2.String())
		require.Equal(t, "168h0m0s", next2.Sub(next.Time).String())

		next = p.Next(start, 10)
		require.Equal(t, "2020-05-11T16:00", next.String())
		require.Equal(t, "1679h0m0s", next.Sub(start.Time).String())
	})
}
