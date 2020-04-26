// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/stretchr/testify/require"
)

func TestUserUnavailable(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)
		// test-user is in PST
		outmd := mustRun(t, SL,
			`/lotto user unavailable --start 2025-01-01T11:00 --finish 2025-01-02T09:30`)
		require.Equal(t, "added unavailable event personal: 2025-01-01T11:00 to 2025-01-02T09:30 to @test-user", outmd.String())

		user := sl.NewUser("")
		err := store.Entity(sl.KeyUser).Load("test-user", user)
		require.NoError(t, err)
		require.Len(t, user.Calendar, 1)
		require.EqualValues(t,
			sl.Unavailable{
				Interval: types.MustParseInterval("2025-01-01T19:00", "2025-01-02T17:30"),
				Reason:   sl.ReasonPersonal,
			},
			*user.Calendar[0])

		mustRunMulti(t, SL, `
				/lotto user unavailable --start 2025-02-01 --finish 2025-02-03
				/lotto user unavailable --start 2025-02-07 --finish 2025-02-10
				/lotto user unavailable --start 2025-06-28 --finish 2025-07-05
			`)

		users := mustRunUsersCalendar(t, SL,
			`/lotto user unavailable --clear --start 2025-01-30T10:00 --finish 2025-02-08T11:00`)
		require.EqualValues(t, []string{"test-user"}, users.TestIDs())
		require.Equal(t, 2, len(users.Get("test-user").Calendar))
		require.EqualValues(t,
			&sl.Unavailable{
				Interval: types.MustParseInterval("2025-01-01T19:00", "2025-01-02T17:30"),
				Reason:   sl.ReasonPersonal,
			},
			users.Get("test-user").Calendar[0])
		require.EqualValues(t,
			&sl.Unavailable{
				Interval: types.MustParseInterval("2025-06-28T07:00", "2025-07-05T07:00"),
				Reason:   sl.ReasonPersonal,
			},
			users.Get("test-user").Calendar[1])
	})
}
