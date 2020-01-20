// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	"testing"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/stretchr/testify/require"
)

const (
	UserIDGuru    = "test-user-guru"
	UserIDServer1 = "test-user-server1"
	UserIDServer2 = "test-user-server2"
	UserIDServer3 = "test-user-server3"
	UserIDWebapp1 = "test-user-webapp1"
	UserIDWebapp2 = "test-user-webapp2"
	UserIDWebapp3 = "test-user-webapp3"
	UserIDMobile1 = "test-user-mobile1"
	UserIDMobile2 = "test-user-mobile2"
)

func UserGuru() *api.User {
	return SkilledUser(UserIDGuru, SkillWebapp, 4, SkillMobile, 4, SkillServer, 4, SkillPlugins, 4)
}
func UserServer1() *api.User {
	return SkilledUser(UserIDServer1, SkillWebapp, 1, SkillServer, 3, SkillPlugins, 1)
}
func UserServer2() *api.User {
	return SkilledUser(UserIDServer2, SkillWebapp, 2, SkillServer, 3, SkillPlugins, 1)
}
func UserServer3() *api.User {
	return SkilledUser(UserIDServer3, SkillWebapp, 1, SkillServer, 3, SkillPlugins, 1)
}
func UserWebapp1() *api.User {
	return SkilledUser(UserIDWebapp1, SkillWebapp, 3, SkillServer, 1)
}
func UserWebapp2() *api.User {
	return SkilledUser(UserIDWebapp2, SkillWebapp, 2, SkillServer, 1)
}
func UserWebapp3() *api.User {
	return SkilledUser(UserIDWebapp3, SkillWebapp, 3, SkillServer, 1)
}
func UserMobile1() *api.User {
	return SkilledUser(UserIDMobile1, SkillWebapp, 1, SkillMobile, 3)
}
func UserMobile2() *api.User {
	return SkilledUser(UserIDMobile2, SkillWebapp, 1, SkillMobile, 3)
}

func AllUsers() api.UserMap {
	return Usermap(
		UserGuru(),
		UserServer1(),
		UserServer2(),
		UserServer3(),
		UserWebapp1(),
		UserWebapp2(),
		UserWebapp3(),
		UserMobile1(),
		UserMobile2(),
	)
}

func User(mattermostUserID string) *api.User {
	return &api.User{
		User: store.NewUser(mattermostUserID),
	}
}

func SkilledUser(mattermostUserID string, skillLevels ...interface{}) *api.User {
	return User(mattermostUserID).WithSkills(Skillmap(skillLevels...))
}

func Usermap(in ...*api.User) api.UserMap {
	users := api.UserMap{}
	for _, u := range in {
		users[u.MattermostUserID] = u.Clone()
	}
	return users
}

func TestUsermap(t *testing.T) {
	require.EqualValues(t,
		api.UserMap{
			UserIDGuru:    UserGuru(),
			UserIDMobile1: UserMobile1(),
		},
		Usermap(
			UserGuru(),
			UserMobile1(),
		),
	)
}

func Skillmap(skillLevels ...interface{}) store.IntMap {
	m := store.IntMap{}
	for i := 0; i < len(skillLevels); i += 2 {
		skill, _ := skillLevels[i].(string)
		level, _ := skillLevels[i+1].(int)
		m[skill] = level
	}
	return m
}

func TestSkillmap(t *testing.T) {
	require.Equal(t,
		store.IntMap{
			"t": 1,
			"z": 2,
		},
		Skillmap("t", 1, "z", 2),
	)
}
