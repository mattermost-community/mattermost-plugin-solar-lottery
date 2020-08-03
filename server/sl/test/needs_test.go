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
		need     sl.Need
		user     *sl.User
		expected bool
	}{
		{
			name:     "Guru server 2",
			need:     C1ServerL2(),
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
			need:     C1ServerL1(),
			user:     UserWebapp1(),
			expected: true,
		},
		{
			name:     "Webapp1 server 2",
			need:     C1ServerL2(),
			user:     UserWebapp1(),
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ok, need := tc.need.QualifyUser(tc.user)
			require.Equal(t, tc.expected, ok)
			if ok {
				require.Equal(t, tc.need.Count()-1, need.Count())
			}
		})
	}
}
func TestUsersQualifiedForNeed(t *testing.T) {
	for _, tc := range []struct {
		name              string
		need              sl.Need
		users             *sl.Users
		expectedQualified *sl.Users
	}{
		{
			name:              "happy server3",
			need:              C1ServerL3(),
			users:             sl.NewUsers(UserGuru(), UserServer1(), UserServer2(), UserServer3(), UserMobile1()),
			expectedQualified: sl.NewUsers(UserGuru(), UserServer1(), UserServer2(), UserServer3()),
		},
		{
			name:              "happy server4",
			need:              C1ServerL4(),
			users:             sl.NewUsers(UserGuru(), UserServer1(), UserServer2(), UserServer3(), UserMobile1()),
			expectedQualified: sl.NewUsers(UserGuru()),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			qualified, need := tc.need.QualifyUsers(tc.users)
			require.Equal(t, tc.expectedQualified, qualified)
			require.Equal(t, tc.need.Count()-int64(tc.expectedQualified.Len()), need.Count())
		})
	}
}
