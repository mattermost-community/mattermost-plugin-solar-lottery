// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func TestHottestRequiredNeed(t *testing.T) {
	needServer1 := store.Need{
		Skill: "server",
		Level: 1,
		Min:   3,
		Max:   -1,
	}

	needServer2 := store.Need{
		Skill: "server",
		Level: 2,
		Min:   2,
		Max:   -1,
	}

	needWebapp2 := store.Need{
		Skill: "webapp",
		Level: 2,
		Min:   1,
		Max:   -1,
	}

	usersServer1 := UserMap{
		testUserGuru.MattermostUserID:    testUserGuru.withWeight(1),
		testUserServer1.MattermostUserID: testUserServer1.withWeight(1),
		testUserServer2.MattermostUserID: testUserServer2.withWeight(1),
		testUserServer3.MattermostUserID: testUserServer3.withWeight(1),
		testUserWebapp1.MattermostUserID: testUserWebapp1.withWeight(1),
		testUserWebapp2.MattermostUserID: testUserWebapp2.withWeight(1),
		testUserWebapp3.MattermostUserID: testUserWebapp3.withWeight(1),
	}

	usersServer2 := UserMap{
		testUserGuru.MattermostUserID:    testUserGuru.withWeight(1),
		testUserServer1.MattermostUserID: testUserServer1.withWeight(1),
		testUserServer2.MattermostUserID: testUserServer2.withWeight(1),
		testUserServer3.MattermostUserID: testUserServer3.withWeight(1),
	}

	usersWebapp2 := UserMap{
		testUserGuru.MattermostUserID:    testUserGuru.withWeight(1),
		testUserServer2.MattermostUserID: testUserServer2.withWeight(1),
		testUserWebapp1.MattermostUserID: testUserWebapp1.withWeight(1),
		testUserWebapp2.MattermostUserID: testUserWebapp2.withWeight(1),
		testUserWebapp3.MattermostUserID: testUserWebapp3.withWeight(1),
	}

	for _, tc := range []struct {
		name          string
		requiredNeeds map[*store.Need]UserMap
		expectedNeed  store.Need
		expectedPool  UserMap
	}{

		{
			name: "happy 2",
			requiredNeeds: map[*store.Need]UserMap{
				&needServer1: usersServer1,
				&needWebapp2: usersWebapp2,
			},
			expectedNeed: needServer1,
			expectedPool: usersServer1,
		},
		{
			name: "happy 3",
			requiredNeeds: map[*store.Need]UserMap{
				&needServer1: usersServer1,
				&needServer2: usersServer2,
				&needWebapp2: usersWebapp2,
			},

			// needServer2 is selected since it has the lowest weight/headcount
			expectedNeed: needServer2,
			expectedPool: usersServer2,
		},
		{
			name: "happy 3 with constraint",
			requiredNeeds: map[*store.Need]UserMap{
				&needServer1:            usersServer1,
				&needServer2:            usersServer2,
				withMax(needWebapp2, 1): usersWebapp2,
			},

			// needServer2 is selected since it has the lowest weight/headcount
			expectedNeed: *withMax(needWebapp2, 1),
			expectedPool: usersWebapp2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			need, pool := hottestRequiredNeed(tc.requiredNeeds)
			require.Equal(t, tc.expectedNeed, *need)
			require.Equal(t, tc.expectedPool, pool)
		})
	}
}

func TestQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name     string
		need     store.Need
		user     *User
		expected bool
	}{
		{
			name:     "Guru server 2",
			need:     store.Need{Skill: "server", Level: 2},
			user:     testUserGuru,
			expected: true,
		},
		{
			name:     "Guru other",
			need:     store.Need{Skill: "other", Level: 2},
			user:     testUserGuru,
			expected: false,
		},
		{
			name:     "Webapp1 server 1",
			need:     store.Need{Skill: "server", Level: 1},
			user:     testUserWebapp1,
			expected: true,
		},
		{
			name:     "Webapp1 server 2",
			need:     store.Need{Skill: "server", Level: 2},
			user:     testUserWebapp1,
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := qualifiedForNeed(tc.user, tc.need)
			require.Equal(t, tc.expected, result)
		})
	}
}
func TestUsersQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name              string
		need              store.Need
		users             UserMap
		expectedQualified UserMap
	}{
		{
			name:              "empty users",
			need:              store.Need{Skill: "server", Level: 2},
			expectedQualified: UserMap{},
		},
		{
			name: "happy server3",
			need: store.Need{Skill: "server", Level: 3},
			users: UserMap{
				testUserGuru.MattermostUserID:    testUserGuru.withWeight(64),
				testUserServer1.MattermostUserID: testUserServer1.withWeight(32),
				testUserServer2.MattermostUserID: testUserServer2.withWeight(32),
				testUserServer3.MattermostUserID: testUserServer3.withWeight(16),
				testUserMobile1.MattermostUserID: testUserMobile1.withWeight(16),
			},
			expectedQualified: UserMap{
				testUserGuru.MattermostUserID:    testUserGuru.withWeight(64),
				testUserServer1.MattermostUserID: testUserServer1.withWeight(32),
				testUserServer2.MattermostUserID: testUserServer2.withWeight(32),
				testUserServer3.MattermostUserID: testUserServer3.withWeight(16),
			},
		},
		{
			name: "happy server4",
			need: store.Need{Skill: "server", Level: 4},
			users: UserMap{
				testUserGuru.MattermostUserID:    testUserGuru.withWeight(64),
				testUserServer1.MattermostUserID: testUserServer1.withWeight(32),
				testUserServer2.MattermostUserID: testUserServer2.withWeight(32),
				testUserServer3.MattermostUserID: testUserServer3.withWeight(16),
				testUserMobile1.MattermostUserID: testUserMobile1.withWeight(16),
			},
			expectedQualified: UserMap{
				testUserGuru.MattermostUserID: testUserGuru.withWeight(64),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			qualified := usersQualifiedForNeed(tc.users, tc.need)
			require.Equal(t, tc.expectedQualified, qualified)
		})
	}
}
func TestUnmetNeeds(t *testing.T) {
	for _, tc := range []struct {
		name          string
		needs         []store.Need
		users         UserMap
		expectedUnmet []store.Need
	}{
		{
			name:          "empty 1",
			needs:         nil,
			users:         nil,
			expectedUnmet: nil,
		},
		{
			name:          "empty 2",
			needs:         []store.Need{},
			users:         UserMap{},
			expectedUnmet: nil,
		},
		{
			name: "happy 1",
			needs: []store.Need{
				store.Need{Skill: "server", Level: 1, Min: 2, Max: -1},
				store.Need{Skill: "webapp", Level: 1, Min: 2, Max: -1},
				store.Need{Skill: "uncovered", Level: 1, Min: 2, Max: -1},
			},
			users: UserMap{
				testUserGuru.MattermostUserID:    testUserGuru,
				testUserMobile1.MattermostUserID: testUserMobile1,
				testUserWebapp1.MattermostUserID: testUserWebapp1,
			},
			expectedUnmet: []store.Need{
				store.Need{Skill: "uncovered", Level: 1, Min: 2, Max: -1},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			unmet := unmetNeeds(tc.needs, tc.users)
			require.Equal(t, tc.expectedUnmet, unmet)
		})
	}
}

func withMax(need store.Need, max int) *store.Need {
	need.Max = max
	return &need
}
