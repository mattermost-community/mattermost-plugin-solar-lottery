// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

// It's a little hacky to have tests here for api functions, but the data is
// nicely set up and it doesn't require in-package visibility...

import (
	"testing"

	"github.com/stretchr/testify/require"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func TestQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name     string
		need     *store.Need
		user     *sl.User
		expected bool
	}{
		{
			name:     "Guru server 2",
			need:     store.NewNeed(SkillServer, 2, 0),
			user:     UserGuru(),
			expected: true,
		},
		{
			name:     "Guru other",
			need:     store.NewNeed("other", 2, 0),
			user:     UserGuru(),
			expected: false,
		},
		{
			name:     "Webapp1 server 1",
			need:     store.NewNeed(SkillServer, 1, 0),
			user:     UserWebapp1(),
			expected: true,
		},
		{
			name:     "Webapp1 server 2",
			need:     store.NewNeed(SkillServer, 2, 0),
			user:     UserWebapp1(),
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := sl.IsUserQualifiedForNeed(tc.user, tc.need)
			require.Equal(t, tc.expected, result)
		})
	}
}
func TestUsersQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name              string
		need              *store.Need
		users             sl.UserMap
		expectedQualified sl.UserMap
	}{
		{
			name:              "empty users",
			need:              store.NewNeed(SkillServer, 2, 0),
			expectedQualified: sl.UserMap{},
		},
		{
			name:              "happy server3",
			need:              store.NewNeed(SkillServer, 3, 0),
			users:             Usermap(UserGuru(), UserServer1(), UserServer2(), UserServer3(), UserMobile1()),
			expectedQualified: Usermap(UserGuru(), UserServer1(), UserServer2(), UserServer3()),
		},
		{
			name:              "happy server4",
			need:              store.NewNeed(SkillServer, 4, 0),
			users:             Usermap(UserGuru(), UserServer1(), UserServer2(), UserServer3(), UserMobile1()),
			expectedQualified: Usermap(UserGuru()),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			qualified := sl.UsersQualifiedForNeed(tc.users, tc.need)
			require.Equal(t, tc.expectedQualified, qualified)
		})
	}
}

func TestUnmetNeeds(t *testing.T) {
	for _, tc := range []struct {
		name          string
		needs         store.Needs
		users         sl.UserMap
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
			users:         sl.UserMap{},
			expectedUnmet: nil,
		},
		{
			name: "happy 1",
			needs: store.Needs{
				store.NewNeed(SkillServer, 1, 2),
				store.NewNeed(SkillWebapp, 1, 2),
				store.NewNeed("uncovered", 1, 2),
			},
			users: Usermap(UserGuru(), UserMobile1(), UserWebapp1()),
			expectedUnmet: store.Needs{
				store.NewNeed("uncovered", 1, 2),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			unmet := sl.UnmetNeeds(tc.needs, tc.users)
			require.Equal(t, tc.expectedUnmet, unmet)
		})
	}
}
