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
	// By testAllUsers has last served of 0
	for _, tc := range []struct {
		name                       string
		size                       int
		needs                      []store.Need
		pool                       UserMap
		chosen                     UserMap
		shiftNumber                int
		shiftStart                 time.Time
		shiftEnd                   time.Time
		expectError                bool
		expectAutofillError        error
		expectCauseCapacity        int
		expectCauseUnmetNeed       *store.Need
		expectCauseUnmetNeeds      []store.Need
		skipSucessResultValidation bool
		expectedChosen             UserMap
		expectedPool               UserMap
	}{
		{
			name:                       "happy",
			size:                       3,
			needs:                      []store.Need{testNeedMobile_L1_Min1, testNeedServer_L1_Min1, testNeedWebapp_L1_Min1},
			pool:                       testAllUsers.Clone(true),
			skipSucessResultValidation: true,
		},
		{
			//
			name:        "happy accepted",
			size:        1,
			shiftNumber: 128,
			needs:       []store.Need{testNeedMobile_L1_Min1},
			pool: usermap(
				testUserMobile1.withLastServed(testRotationID, 1),
				testUserMobile2.withLastServed(testRotationID, 127),
			),
			expectedChosen: usermap(testUserMobile1.withLastServed(testRotationID, 1)),
			expectedPool:   usermap(testUserMobile2.withLastServed(testRotationID, 127)),
		},
		{
			name:                  "ErrInsufficientForNeeds",
			size:                  1,
			needs:                 []store.Need{testNeedMobile_L1_Min1},
			pool:                  usermap(),
			chosen:                usermap(testUserMobile2),
			expectAutofillError:   ErrInsufficientForNeeds,
			expectCauseUnmetNeeds: []store.Need{testNeedMobile_L1_Min1},
			expectCauseUnmetNeed:  &testNeedMobile_L1_Min1,
		},
		{
			name:                  "ErrInsufficientForNeeds in acceptUser",
			size:                  1,
			needs:                 []store.Need{testNeedMobile_L1_Min2},
			pool:                  usermap(testUserMobile1),
			chosen:                usermap(testUserMobile2),
			expectAutofillError:   ErrInsufficientForNeeds,
			expectCauseUnmetNeeds: []store.Need{testNeedMobile_L1_Min1},
			expectCauseUnmetNeed:  &testNeedMobile_L1_Min1,
		},
		{
			name:                "ErrInsufficientForSize",
			size:                2,
			pool:                usermap(),
			chosen:              usermap(testUserMobile2),
			expectAutofillError: ErrInsufficientForSize,
			expectCauseCapacity: 1,
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
				tc.shiftStart,
				tc.shiftEnd,
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
			require.Equal(t, tc.expectCauseCapacity, afErr.causeCapacity)
			require.Equal(t, tc.expectCauseUnmetNeeds, afErr.causeUnmetNeeds)
			require.Equal(t, tc.expectCauseUnmetNeed, afErr.causeUnmetNeed)
		})
	}
}
