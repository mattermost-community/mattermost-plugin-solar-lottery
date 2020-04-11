// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/test"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

var (
	longTimeAgo = types.MustParseTime("2010-01-01")
	recently    = types.MustParseTime("2020-01-01")
)

func TestFill(t *testing.T) {
	// This test makes a user get picked by inflating the probability of her
	// being selected. This is accomplished by setting the LastServed to 0.
	// Conversely, setting it to 127 for shift 128 makes it extremely unlikely
	// by comparison. By default, testAllUsers has LastServed of 0.
	for _, tc := range []struct {
		name                       string
		require                    *sl.Needs
		limit                      *sl.Needs
		pool                       *sl.Users
		assigned                   *sl.Users
		time                       types.Time
		expectError                bool
		expectFillError            error
		expectFailedNeed           sl.Need
		expectUnmetNeeds           *sl.Needs
		skipSucessResultValidation bool
		expectFilled               *sl.Users
		expectPool                 *sl.Users
	}{
		{
			name:                       "happy",
			require:                    sl.NewNeeds(test.C1_Mobile_L1(), test.C1_Server_L1(), test.C1_Webapp_L1()),
			pool:                       test.AllUsers(),
			skipSucessResultValidation: true,
		},
		{
			name:    "happy accepted",
			require: sl.NewNeeds(test.C1_Mobile_L1()),
			pool: sl.NewUsers(
				test.UserMobile1().WithLastServed(test.RotationID, longTimeAgo),
				test.UserMobile2().WithLastServed(test.RotationID, recently),
			),
			expectFilled: sl.NewUsers(test.UserMobile1().WithLastServed(test.RotationID, longTimeAgo)),
			expectPool:   sl.NewUsers(test.UserMobile2().WithLastServed(test.RotationID, recently)),
		},
		{
			name: "accepted another maxed out",
			require: sl.NewNeeds(
				test.C1_Mobile_L1(),
				test.C1_Server_L1(),
			),
			limit: sl.NewNeeds(
				test.C1_Server_L1(),
			),
			assigned: sl.NewUsers(
				test.UserServer1(), // meets and maxes out the server need, leaving only mobile
			),
			pool: sl.NewUsers(
				test.UserGuru().WithLastServed(test.RotationID, longTimeAgo), // higher probability for Mobile, but should be rejected on limit
				test.UserMobile2().WithLastServed(test.RotationID, recently),
			),
			expectFilled: sl.NewUsers(
				test.UserMobile2().WithLastServed(test.RotationID, recently),
			),
			expectPool: sl.NewUsers(),
		},
		{
			name:             "Err Insufficient simple",
			require:          sl.NewNeeds(test.C2_Mobile_L1()),
			pool:             sl.NewUsers(test.UserServer1(), test.UserMobile1()),
			expectFillError:  sl.ErrFillInsufficient,
			expectUnmetNeeds: sl.NewNeeds(test.C1_Mobile_L1()),
			expectFailedNeed: test.C1_Mobile_L1(),
		},
		{
			name:             "Err Insufficient with assigned",
			require:          sl.NewNeeds(test.C2_Mobile_L1()),
			assigned:         sl.NewUsers(test.UserMobile1()),
			pool:             sl.NewUsers(test.UserServer1(), test.UserServer2()),
			expectFillError:  sl.ErrFillInsufficient,
			expectUnmetNeeds: sl.NewNeeds(test.C1_Mobile_L1()),
			expectFailedNeed: test.C1_Mobile_L1(),
		},
		{
			name:             "Err Limit",
			limit:            sl.NewNeeds(test.C1_Any()),
			require:          sl.NewNeeds(test.C2_Mobile_L1()),
			assigned:         sl.NewUsers(test.UserMobile1()),
			pool:             sl.NewUsers(test.UserMobile2(), test.UserGuru()),
			expectFillError:  sl.ErrFillInsufficient,
			expectUnmetNeeds: sl.NewNeeds(test.C1_Mobile_L1()),
			expectFailedNeed: test.C1_Mobile_L1(),
		},
		{
			name:         "Noop on preassigned",
			limit:        sl.NewNeeds(test.C2_Any()),
			require:      sl.NewNeeds(test.C2_Mobile_L1()),
			assigned:     sl.NewUsers(test.UserMobile1(), test.UserMobile2()),
			pool:         sl.NewUsers(test.UserWebapp2(), test.UserGuru()),
			expectPool:   sl.NewUsers(test.UserWebapp2(), test.UserGuru()),
			expectFilled: sl.NewUsers(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := makeTestFiller(t, tc.pool, tc.assigned, tc.require, tc.limit)
			filled, err := f.fill()
			if !tc.expectError && tc.expectFillError == nil {
				require.NoError(t, err)
				if tc.skipSucessResultValidation {
					return
				}

				require.EqualValues(t, tc.expectFilled.TestArray(), filled.TestArray())
				require.EqualValues(t, tc.expectPool.TestArray(), f.pool.TestArray())
				return
			}

			require.Error(t, err)
			if tc.expectFillError == nil {
				return
			}

			ferr, _ := err.(*sl.FillError)
			require.NotNil(t, ferr)
			require.Equal(t, tc.expectFillError, ferr.Err)
			require.Equal(t, tc.expectUnmetNeeds.AsArray(), ferr.UnmetNeeds.AsArray())
			require.EqualValues(t, &tc.expectFailedNeed, ferr.FailedNeed)
		})
	}
}

func makeTestFiller(t testing.TB, pool, assigned *sl.Users, require, limit *sl.Needs) *fill {
	if pool.IsEmpty() {
		pool = sl.NewUsers()
	}
	r := sl.NewRotation()
	r.RotationID = test.RotationID
	r.Beginning = types.MustParseTime("2020-01-01")
	r.Users = pool

	task := sl.NewTask(r.RotationID)
	task.TaskID = r.RotationID + "#2020-02-02"
	if assigned == nil {
		task.Users = sl.NewUsers()
	} else {
		task.Users = assigned
	}

	if !limit.IsEmpty() {
		task.Limit = limit
	}
	if !require.IsEmpty() {
		task.Require = require
	}
	return newFill(r, task, types.MustParseTime("2020-03-01"),
		// &bot.TestLogger{TB: t},
		&bot.NilLogger{},
	)
}
