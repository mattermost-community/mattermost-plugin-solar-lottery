// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

// A rotation requires 3 users, 4 different skills.

var testLastServed0 = store.IntMap{testRotationID: 0}
var testLastServed1 = store.IntMap{testRotationID: 1}
var testLastServed2 = store.IntMap{testRotationID: 2}

var testUserGuru = newUser(&store.User{
	MattermostUserID: "test-user-guru",
	SkillLevels: store.IntMap{
		testSkillWebapp:  4,
		testSkillMobile:  4,
		testSkillServer:  4,
		testSkillPlugins: 4,
	},
	LastServed: testLastServed0,
})

var testUserServer1 = newUser(&store.User{
	MattermostUserID: "test-user-server1",
	SkillLevels: store.IntMap{
		testSkillWebapp:  1,
		testSkillServer:  3,
		testSkillPlugins: 1,
	},
	LastServed: testLastServed0,
})

var testUserServer2 = newUser(&store.User{
	MattermostUserID: "test-user-server2",
	SkillLevels: store.IntMap{
		testSkillWebapp:  2,
		testSkillServer:  3,
		testSkillPlugins: 1,
	},
	LastServed: testLastServed0,
})

var testUserServer3 = newUser(&store.User{
	MattermostUserID: "test-user-server3",
	SkillLevels: store.IntMap{
		testSkillWebapp:  1,
		testSkillServer:  3,
		testSkillPlugins: 1,
	},
	LastServed: testLastServed0,
})

var testUserWebapp1 = newUser(&store.User{
	MattermostUserID: "test-user-webapp1",
	SkillLevels: store.IntMap{
		testSkillWebapp: 3,
		testSkillServer: 1,
	},
	LastServed: testLastServed0,
})

var testUserWebapp2 = newUser(&store.User{
	MattermostUserID: "test-user-webapp2",
	SkillLevels: store.IntMap{
		testSkillWebapp: 2,
		testSkillServer: 1,
	},
	LastServed: testLastServed0,
})

var testUserWebapp3 = newUser(&store.User{
	MattermostUserID: "test-user-webapp3",
	SkillLevels: store.IntMap{
		testSkillWebapp: 3,
		testSkillServer: 1,
	},
	LastServed: testLastServed0,
})

var testUserMobile1 = newUser(&store.User{
	MattermostUserID: "test-user-mobile1",
	SkillLevels: store.IntMap{
		testSkillWebapp: 1,
		testSkillMobile: 3,
	},
	LastServed: testLastServed0,
})

var testUserMobile2 = newUser(&store.User{
	MattermostUserID: "test-user-mobile2",
	SkillLevels: store.IntMap{
		testSkillWebapp: 1,
		testSkillMobile: 3,
	},
	LastServed: testLastServed0,
})

var testAllUsers = usermap(
	testUserGuru,
	testUserServer1,
	testUserServer2,
	testUserServer3,
	testUserWebapp1,
	testUserWebapp2,
	testUserWebapp3,
	testUserMobile1,
	testUserMobile2,
)

func newUser(su *store.User) *User {
	return &User{
		User: su,
	}
}

func usermap(in ...*User) UserMap {
	users := UserMap{}
	for _, u := range in {
		u = u.Clone()
		users[u.MattermostUserID] = u.Clone()
	}
	return users
}

func (user *User) withWeight(weight float64) *User {
	newUser := user.Clone()
	newUser.weight = weight
	return newUser
}

func (user *User) withLastServed(rotationID string, shiftNumber int) *User {
	newUser := user.Clone()
	newUser.LastServed[rotationID] = shiftNumber
	return newUser
}
