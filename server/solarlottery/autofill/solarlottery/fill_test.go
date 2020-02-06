// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/autofill"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery/test"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func TestFillOneAllHappy(t *testing.T) {
	// This test makes a user get picked by inflating the probability of her
	// being selected. This is accomplished by setting the LastServed to 0.
	// Conversely, setting it to 127 for shift 128 makes it extremely unlikely
	// by comparison. By default, testAllUsers has LastServed of 0.
	for _, tc := range []struct {
		name                       string
		size                       int
		needs                      store.Needs
		pool                       sl.UserMap
		chosen                     sl.UserMap
		shiftNumber                int
		expectError                bool
		expectAutofillError        error
		expectUnmetCapacity        int
		expectUnmetNeed            *store.Need
		expectUnmetNeeds           store.Needs
		skipSucessResultValidation bool
		expectedChosen             sl.UserMap
		expectedPool               sl.UserMap
	}{
		{
			name: "happy",
			size: 3,
			needs: store.Needs{
				test.NeedMobile_L1_Min1(),
				test.NeedServer_L1_Min1(),
				test.NeedWebapp_L1_Min1()},
			pool:                       test.AllUsers(),
			skipSucessResultValidation: true,
		},
		{
			name:        "happy accepted",
			size:        1,
			shiftNumber: 128,
			needs:       store.Needs{test.NeedMobile_L1_Min1()},
			pool: test.Usermap(
				test.UserMobile1().WithLastServed(test.RotationID, 1),
				test.UserMobile2().WithLastServed(test.RotationID, 127),
			),
			expectedChosen: test.Usermap(test.UserMobile1().WithLastServed(test.RotationID, 1)),
			expectedPool:   test.Usermap(test.UserMobile2().WithLastServed(test.RotationID, 127)),
		},
		{
			name:        "accepted another maxed out",
			size:        1,
			shiftNumber: 128,
			needs: store.Needs{
				test.NeedMobile_L1_Min1(),
				test.NeedServer_L1_Min1().WithMax(1),
			},
			chosen: test.Usermap(
				test.UserServer1(), // meets and maxes out the server need
			),
			pool: test.Usermap(
				test.UserGuru(), // LastServed 0, higher probability for Mobile
				test.UserMobile2().WithLastServed(test.RotationID, 127),
			),
			expectedChosen: test.Usermap(
				test.UserServer1(),
				test.UserMobile2().WithLastServed(test.RotationID, 127),
			),
			expectedPool: sl.UserMap{},
		},
		{
			name:                "ErrInsufficientForNeeds",
			size:                1,
			needs:               store.Needs{test.NeedMobile_L1_Min2()},
			pool:                test.Usermap(),
			chosen:              test.Usermap(test.UserMobile2()),
			expectAutofillError: autofill.ErrInsufficientForNeeds,
			expectUnmetNeeds:    store.Needs{test.NeedMobile_L1_Min1()},
			expectUnmetNeed:     test.NeedMobile_L1_Min1(),
			expectUnmetCapacity: 1,
		},
		{
			name:                "ErrInsufficientForNeeds in acceptUser",
			size:                1,
			needs:               store.Needs{test.NeedMobile_L1_Min3()},
			pool:                test.Usermap(test.UserMobile1()),
			chosen:              test.Usermap(test.UserMobile2()),
			expectAutofillError: autofill.ErrInsufficientForNeeds,
			expectUnmetNeeds:    store.Needs{test.NeedMobile_L1_Min1()},
			expectUnmetNeed:     test.NeedMobile_L1_Min1(),
		},
		{
			name:                "ErrInsufficientForSize",
			size:                2,
			pool:                test.Usermap(),
			chosen:              test.Usermap(test.UserMobile2()),
			expectAutofillError: autofill.ErrInsufficientForSize,
			expectUnmetCapacity: 1,
		},
		{
			name:                       "does not enforce ErrSizeExceeded",
			size:                       1,
			pool:                       test.Usermap(test.UserMobile1()),
			chosen:                     test.Usermap(test.UserMobile2()),
			skipSucessResultValidation: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			af, err := makeTestAutofill(t, tc.size, tc.needs, tc.pool, tc.chosen, tc.shiftNumber)
			require.NoError(t, err)

			err = af.fillOne()
			if !tc.expectError && tc.expectAutofillError == nil {
				require.NoError(t, err)
				if tc.skipSucessResultValidation {
					return
				}

				require.EqualValues(t, tc.expectedChosen.IDMap(), af.chosen.IDMap(), "chosen")
				require.EqualValues(t, tc.expectedPool.IDMap(), af.pool.IDMap(), "pool")
				return
			}

			require.Error(t, err)
			if tc.expectAutofillError == nil {
				return
			}

			afErr, _ := err.(*autofill.Error)
			require.NotNil(t, afErr)
			require.Equal(t, tc.expectAutofillError, afErr.Err)
			require.Equal(t, tc.expectUnmetCapacity, afErr.UnmetCapacity)
			require.Equal(t, tc.expectUnmetNeeds, afErr.UnmetNeeds)
			require.Equal(t, tc.expectUnmetNeed, afErr.UnmetNeed)
		})
	}
}

func TestMeetsConstraints(t *testing.T) {
	// By testAllUsers has last served of 0
	for _, tc := range []struct {
		name             string
		user             *sl.User
		constrainedNeeds store.Needs
		expectedResult   bool
	}{
		{
			name:             "meets 1",
			constrainedNeeds: store.Needs{test.NeedMobile_L1_Min1().WithMax(1)},
			user:             test.UserServer1(),
			expectedResult:   true,
		},
		{
			name: "meets 2",
			constrainedNeeds: store.Needs{
				test.NeedMobile_L1_Min1().WithMax(5),
				test.NeedServer_L1_Min1().WithMax(1),
			},
			user:           test.UserServer1(),
			expectedResult: true,
		},
		{
			name:           "meets 3",
			user:           test.UserServer1(),
			expectedResult: true,
		},
		{
			name: "doesnt meet",
			constrainedNeeds: store.Needs{
				test.NeedMobile_L1_Min1().WithMax(5),
				test.NeedServer_L2_Min1().WithMax(1),
				test.NeedServer_L1_Min1().WithMax(0),
			},
			user:           test.UserServer1(),
			expectedResult: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			af, err := makeTestAutofill(t, 10, nil, nil, nil, 0)
			require.NoError(t, err)
			af.constrainedNeeds = tc.constrainedNeeds

			meets := af.meetsConstraints(tc.user)
			require.Equal(t, tc.expectedResult, meets)
		})
	}
}

func makeTestAutofill(t testing.TB, size int, needs store.Needs,
	pool sl.UserMap, chosen sl.UserMap, shiftNumber int) (*fill, error) {
	return makeAutofill(test.RotationID, size, needs, pool, chosen, shiftNumber, time.Time{}, time.Time{},
		// &bot.TestLogger{TB: t},
		&bot.NilLogger{},
	)
}
