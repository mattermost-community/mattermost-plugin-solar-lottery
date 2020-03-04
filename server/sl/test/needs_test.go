// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

// It's a little hacky to have tests here for api functions, but the data is
// nicely set up and it doesn't require in-package visibility...

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
)

func TestQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name     string
		need     *sl.Need
		user     *sl.User
		expected bool
	}{
		{
			name:     "Guru server 2",
			need:     C1_Server_L2(),
			user:     UserGuru(),
			expected: true,
		},
		{
			name:     "Guru other",
			need:     sl.NewNeed(0, sl.NewSkillLevel("other", 2)),
			user:     UserGuru(),
			expected: false,
		},
		{
			name:     "Webapp1 server 1",
			need:     C1_Server_L1(),
			user:     UserWebapp1(),
			expected: true,
		},
		{
			name:     "Webapp1 server 2",
			need:     C1_Server_L2(),
			user:     UserWebapp1(),
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.user.IsQualified(tc.need.SkillLevel)
			require.Equal(t, tc.expected, result)
		})
	}
}
func TestUsersQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name              string
		need              *sl.Need
		users             *sl.Users
		expectedQualified *sl.Users
	}{
		{
			name:              "happy server3",
			need:              C1_Server_L3(),
			users:             sl.NewUsers(UserGuru(), UserServer1(), UserServer2(), UserServer3(), UserMobile1()),
			expectedQualified: sl.NewUsers(UserGuru(), UserServer1(), UserServer2(), UserServer3()),
		},
		{
			name:              "happy server4",
			need:              C1_Server_L4(),
			users:             sl.NewUsers(UserGuru(), UserServer1(), UserServer2(), UserServer3(), UserMobile1()),
			expectedQualified: sl.NewUsers(UserGuru()),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			qualified := tc.users.Qualified(tc.need.SkillLevel)
			require.Equal(t, tc.expectedQualified, qualified)
		})
	}
}

// func TestUnmetNeeds(t *testing.T) {
// 	for _, tc := range []struct {
// 		name          string
// 		needs         []*sl.Need
// 		users         sl.Users
// 		expectedUnmet []*sl.Need
// 	}{
// 		{
// 			name:          "empty 1",
// 			needs:         nil,
// 			users:         nil,
// 			expectedUnmet: nil,
// 		},
// 		{
// 			name:          "empty 2",
// 			needs:         []sl.Need{},
// 			users:         sl.NewUsers(),
// 			expectedUnmet: nil,
// 		},
// 		{
// 			name: "happy 1",
// 			needs: []*sl.Need{
// 				sl.NewNeed(SkillServer, 1, 2),
// 				sl.NewNeed(SkillWebapp, 1, 2),
// 				sl.NewNeed("uncovered", 1, 2),
// 			},
// 			users: Usermap(UserGuru(), UserMobile1(), UserWebapp1()),
// 			expectedUnmet: []*sl.Need{
// 				sl.NewNeed("uncovered", 1, 2),
// 			},
// 		},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			unmet := sl.UnmetNeeds(tc.needs, tc.users)
// 			require.Equal(t, tc.expectedUnmet, unmet)
// 		})
// 	}
// }
