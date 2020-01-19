// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
		needs                      []*store.Need
		pool                       UserMap
		chosen                     UserMap
		shiftNumber                int
		expectError                bool
		expectAutofillError        error
		expectUnmetCapacity        int
		expectUnmetNeed            *store.Need
		expectUnmetNeeds           []*store.Need
		skipSucessResultValidation bool
		expectedChosen             UserMap
		expectedPool               UserMap
	}{
		{
			name:                       "happy",
			size:                       3,
			needs:                      []*store.Need{testNeedMobile_L1_Min1(), testNeedServer_L1_Min1(), testNeedWebapp_L1_Min1()},
			pool:                       testAllUsers.Clone(true),
			skipSucessResultValidation: true,
		},
		{
			name:        "happy accepted",
			size:        1,
			shiftNumber: 128,
			needs:       []*store.Need{testNeedMobile_L1_Min1()},
			pool: usermap(
				testUserMobile1.withLastServed(testRotationID, 1),
				testUserMobile2.withLastServed(testRotationID, 127),
			),
			expectedChosen: usermap(testUserMobile1.withLastServed(testRotationID, 1)),
			expectedPool:   usermap(testUserMobile2.withLastServed(testRotationID, 127)),
		},
		{
			name:        "accepted another maxed out",
			size:        1,
			shiftNumber: 128,
			needs: []*store.Need{
				testNeedMobile_L1_Min1(),
				testNeedServer_L1_Min1().WithMax(1),
			},
			chosen: usermap(
				testUserServer1, // meets and maxes out the server need
			),
			pool: usermap(
				testUserGuru, // LastServed 0, higher probability for Mobile
				testUserMobile2.withLastServed(testRotationID, 127),
			),
			expectedChosen: usermap(
				testUserServer1,
				testUserMobile2.withLastServed(testRotationID, 127),
			),
			expectedPool: UserMap{},
		},
		{
			name:                "ErrInsufficientForNeeds",
			size:                1,
			needs:               []*store.Need{testNeedMobile_L1_Min2()},
			pool:                usermap(),
			chosen:              usermap(testUserMobile2),
			expectAutofillError: ErrInsufficientForNeeds,
			expectUnmetNeeds:    []*store.Need{testNeedMobile_L1_Min1()},
			expectUnmetNeed:     testNeedMobile_L1_Min1(),
			expectUnmetCapacity: 1,
		},
		{
			name:                "ErrInsufficientForNeeds in acceptUser",
			size:                1,
			needs:               []*store.Need{testNeedMobile_L1_Min3()},
			pool:                usermap(testUserMobile1),
			chosen:              usermap(testUserMobile2),
			expectAutofillError: ErrInsufficientForNeeds,
			expectUnmetNeeds:    []*store.Need{testNeedMobile_L1_Min1()},
			expectUnmetNeed:     testNeedMobile_L1_Min1(),
		},
		{
			name:                "ErrInsufficientForSize",
			size:                2,
			pool:                usermap(),
			chosen:              usermap(testUserMobile2),
			expectAutofillError: ErrInsufficientForSize,
			expectUnmetCapacity: 1,
		},
		{
			name:                       "does not enforce ErrSizeExceeded",
			size:                       1,
			pool:                       usermap(testUserMobile1),
			chosen:                     usermap(testUserMobile2),
			skipSucessResultValidation: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := &api{
				Logger: &bot.NilLogger{},
			}
			af, err := api.makeAutofillImpl(
				testRotationID,
				tc.size,
				tc.needs,
				tc.pool,
				tc.chosen,
				tc.shiftNumber,
				time.Time{},
				time.Time{},
			)
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

			afErr, _ := err.(*autofillError)
			require.NotNil(t, afErr)
			require.Equal(t, tc.expectAutofillError, afErr.orig)
			require.Equal(t, tc.expectUnmetCapacity, afErr.unmetCapacity)
			require.Equal(t, tc.expectUnmetNeeds, afErr.unmetNeeds)
			require.Equal(t, tc.expectUnmetNeed, afErr.unmetNeed)
		})
	}
}

func TestMeetsConstraints(t *testing.T) {
	// By testAllUsers has last served of 0
	for _, tc := range []struct {
		name             string
		user             *User
		constrainedNeeds []*store.Need
		expectedResult   bool
	}{
		{
			name:             "meets 1",
			constrainedNeeds: []*store.Need{testNeedMobile_L1_Min1().WithMax(1)},
			user:             testUserServer1,
			expectedResult:   true,
		},
		{
			name: "meets 2",
			constrainedNeeds: []*store.Need{
				testNeedMobile_L1_Min1().WithMax(5),
				testNeedServer_L1_Min1().WithMax(1),
			},
			user:           testUserServer1,
			expectedResult: true,
		},
		{
			name:           "meets 3",
			user:           testUserServer1,
			expectedResult: true,
		},
		{
			name: "doesnt meet",
			constrainedNeeds: []*store.Need{
				testNeedMobile_L1_Min1().WithMax(5),
				testNeedServer_L2_Min1().WithMax(1),
				testNeedServer_L1_Min1().WithMax(0),
			},
			user:           testUserServer1,
			expectedResult: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := &api{
				Logger: &bot.NilLogger{},
			}
			af, err := api.makeAutofillImpl(
				testRotationID,
				10,
				nil,
				nil,
				nil,
				0,
				time.Time{},
				time.Time{},
			)
			require.NoError(t, err)
			af.constrainedNeeds = tc.constrainedNeeds

			meets := af.meetsConstraints(tc.user)
			require.Equal(t, tc.expectedResult, meets)
		})
	}
}
