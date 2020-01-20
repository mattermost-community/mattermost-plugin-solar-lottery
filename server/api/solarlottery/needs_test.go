// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

// It's a little hacky to have tests here for api functions, but the data is
// nicely set up and it doesn't require in-package visibility...

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/test"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func TestHottestRequiredNeed(t *testing.T) {
	usersServer1 := test.Usermap(
		test.UserGuru(),
		test.UserServer1(),
		test.UserServer2(),
		test.UserServer3(),
		test.UserWebapp1(),
		test.UserWebapp2(),
		test.UserWebapp3(),
	)
	usersServer2 := test.Usermap(
		test.UserGuru(),
		test.UserServer1(),
		test.UserServer2(),
		test.UserServer3(),
	)
	usersWebapp2 := test.Usermap(
		test.UserGuru(),
		test.UserServer2(),
		test.UserWebapp1(),
		test.UserWebapp2(),
		test.UserWebapp3(),
	)

	for _, tc := range []struct {
		name          string
		requiredNeeds store.Needs
		needPools     map[string]api.UserMap
		expectedNeed  *store.Need
		expectedPool  api.UserMap
	}{
		{
			name: "happy 2",
			requiredNeeds: store.Needs{
				test.NeedServer_L1_Min3(),
				test.NeedWebapp_L2_Min1(),
			},
			needPools: map[string]api.UserMap{
				test.NeedServer_L1_Min3().SkillLevel(): usersServer1,
				test.NeedWebapp_L2_Min1().SkillLevel(): usersWebapp2,
			},
			expectedNeed: test.NeedServer_L1_Min3(),
			expectedPool: usersServer1,
		},
		{
			name: "happy 3",
			requiredNeeds: store.Needs{
				test.NeedServer_L1_Min3(),
				test.NeedServer_L2_Min2(),
				test.NeedWebapp_L2_Min1(),
			},
			needPools: map[string]api.UserMap{
				test.NeedServer_L1_Min3().SkillLevel(): usersServer1,
				test.NeedServer_L2_Min2().SkillLevel(): usersServer2,
				test.NeedWebapp_L2_Min1().SkillLevel(): usersWebapp2,
			},

			// testNeedServer_L2_Min2 is selected since it has the lowest weight/headcount
			expectedNeed: test.NeedServer_L2_Min2(),
			expectedPool: usersServer2,
		},
		{
			name: "happy 3 with constraint",
			requiredNeeds: store.Needs{
				test.NeedServer_L1_Min3(),
				test.NeedServer_L2_Min2(),
				test.NeedWebapp_L2_Min1().WithMax(1),
			},
			needPools: map[string]api.UserMap{
				test.NeedServer_L1_Min3().SkillLevel(): usersServer1,
				test.NeedServer_L2_Min2().SkillLevel(): usersServer2,
				test.NeedWebapp_L2_Min1().SkillLevel(): usersWebapp2,
			},

			// testNeedServer_L2_Min2 is selected since it has the lowest weight/headcount
			expectedNeed: test.NeedWebapp_L2_Min1().WithMax(1),
			expectedPool: usersWebapp2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			af, err := makeTestAutofill(t, 10, nil, nil, nil, 0)
			require.NoError(t, err)
			need, pool := af.hottestRequiredNeed(tc.requiredNeeds, tc.needPools)
			require.Equal(t, tc.expectedNeed, need)
			require.Equal(t, tc.expectedPool, pool)
		})
	}
}

func TestQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name     string
		need     *store.Need
		user     *api.User
		expected bool
	}{
		{
			name:     "Guru server 2",
			need:     store.NewNeed(test.SkillServer, 2, 0),
			user:     test.UserGuru(),
			expected: true,
		},
		{
			name:     "Guru other",
			need:     store.NewNeed("other", 2, 0),
			user:     test.UserGuru(),
			expected: false,
		},
		{
			name:     "Webapp1 server 1",
			need:     store.NewNeed(test.SkillServer, 1, 0),
			user:     test.UserWebapp1(),
			expected: true,
		},
		{
			name:     "Webapp1 server 2",
			need:     store.NewNeed(test.SkillServer, 2, 0),
			user:     test.UserWebapp1(),
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := api.IsUserQualifiedForNeed(tc.user, tc.need)
			require.Equal(t, tc.expected, result)
		})
	}
}
func TestUsersQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name              string
		need              *store.Need
		users             api.UserMap
		expectedQualified api.UserMap
	}{
		{
			name:              "empty users",
			need:              store.NewNeed(test.SkillServer, 2, 0),
			expectedQualified: api.UserMap{},
		},
		{
			name: "happy server3",
			need: store.NewNeed(test.SkillServer, 3, 0),
			users: test.Usermap(
				test.UserGuru(),
				test.UserServer1(),
				test.UserServer2(),
				test.UserServer3(),
				test.UserMobile1(),
			),
			expectedQualified: test.Usermap(
				test.UserGuru(),
				test.UserServer1(),
				test.UserServer2(),
				test.UserServer3(),
			),
		},
		{
			name: "happy server4",
			need: store.NewNeed(test.SkillServer, 4, 0),
			users: test.Usermap(
				test.UserGuru(),
				test.UserServer1(),
				test.UserServer2(),
				test.UserServer3(),
				test.UserMobile1()),
			expectedQualified: test.Usermap(test.UserGuru()),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			qualified := api.UsersQualifiedForNeed(tc.users, tc.need)
			require.Equal(t, tc.expectedQualified, qualified)
		})
	}
}

func TestUnmetNeeds(t *testing.T) {
	for _, tc := range []struct {
		name          string
		needs         store.Needs
		users         api.UserMap
		expectedUnmet store.Needs
	}{
		{
			name:          "empty 1",
			needs:         nil,
			users:         nil,
			expectedUnmet: nil,
		},
		{
			name:          "empty 2",
			needs:         store.Needs{},
			users:         api.UserMap{},
			expectedUnmet: nil,
		},
		{
			name: "happy 1",
			needs: store.Needs{
				store.NewNeed(test.SkillServer, 1, 2),
				store.NewNeed(test.SkillWebapp, 1, 2),
				store.NewNeed("uncovered", 1, 2),
			},
			users: test.Usermap(test.UserGuru(), test.UserMobile1(), test.UserWebapp1()),
			expectedUnmet: store.Needs{
				store.NewNeed("uncovered", 1, 2),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			unmet := api.UnmetNeeds(tc.needs, tc.users)
			require.Equal(t, tc.expectedUnmet, unmet)
		})
	}
}
