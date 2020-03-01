// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
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

func UserGuru() *sl.User {
	return SkilledUser(UserIDGuru, Webapp, 4, Mobile, 4, Server, 4, Plugins, 4)
}
func UserServer1() *sl.User {
	return SkilledUser(UserIDServer1, Webapp, 1, Server, 3, Plugins, 1)
}
func UserServer2() *sl.User {
	return SkilledUser(UserIDServer2, Webapp, 2, Server, 3, Plugins, 1)
}
func UserServer3() *sl.User {
	return SkilledUser(UserIDServer3, Webapp, 1, Server, 3, Plugins, 1)
}
func UserWebapp1() *sl.User {
	return SkilledUser(UserIDWebapp1, Webapp, 3, Server, 1)
}
func UserWebapp2() *sl.User {
	return SkilledUser(UserIDWebapp2, Webapp, 2, Server, 1)
}
func UserWebapp3() *sl.User {
	return SkilledUser(UserIDWebapp3, Webapp, 3, Server, 1)
}
func UserMobile1() *sl.User {
	return SkilledUser(UserIDMobile1, Webapp, 1, Mobile, 3)
}
func UserMobile2() *sl.User {
	return SkilledUser(UserIDMobile2, Webapp, 1, Mobile, 3)
}

func AllUsers() sl.Users {
	return sl.NewUsers(
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

func SkilledUser(mattermostUserID types.ID, skillLevels ...interface{}) *sl.User {
	return sl.NewUser(mattermostUserID).WithSkills(Skillmap(skillLevels...))
}

func Skillmap(skillLevels ...interface{}) *types.IntIndex {
	m := types.NewIntIndex()
	for i := 0; i < len(skillLevels); i += 2 {
		skill, _ := skillLevels[i].(string)
		level, _ := skillLevels[i+1].(int)
		m.Set(types.ID(skill), int64(level))
	}
	return m
}
