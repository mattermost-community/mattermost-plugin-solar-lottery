// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Users struct {
	types.ValueSet // of *User
}

func NewUsers(uu ...*User) *Users {
	users := Users{
		ValueSet: *types.NewValueSet(&usersArray{}),
	}
	for _, user := range uu {
		users.Set(user)
	}
	return &users
}

func (users Users) Get(id types.ID) *User {
	return users.ValueSet.Get(id).(*User)
}

func (users Users) MarkdownWithSkills() string {
	out := []string{}
	for _, user := range users.AsArray() {
		out = append(out, fmt.Sprintf("%s %s", user.Markdown(), user.MarkdownSkills()))
	}
	return strings.Join(out, ", ")
}

func (users Users) Markdown() string {
	out := []string{}
	for _, user := range users.AsArray() {
		out = append(out, user.Markdown())
	}
	return strings.Join(out, ", ")
}

func (users Users) String() string {
	out := []string{}
	for _, user := range users.AsArray() {
		out = append(out, user.String())
	}
	return strings.Join(out, ", ")
}

// TestArray returns all users, sorted by MattermostUserID. It is used in
// testing, so it returns []User rather than a []*User to make it easier to
// compare with expected results.
func (users Users) TestArray() []User {
	out := []User{}
	for _, id := range users.TestIDs() {
		user := users.Get(types.ID(id))
		out = append(out, *user)
	}
	return out
}

func (users Users) AsArray() []*User {
	a := usersArray{}
	users.ValueSet.AsArray(&a)
	return []*User(a)
}

type usersArray []*User

func (p usersArray) Len() int                   { return len(p) }
func (p usersArray) GetAt(n int) types.Value    { return p[n] }
func (p usersArray) SetAt(n int, v types.Value) { p[n] = v.(*User) }

func (p usersArray) InstanceOf() types.ValueArray {
	inst := make(usersArray, 0)
	return &inst
}
func (p *usersArray) Ref() interface{} { return &p }
func (p *usersArray) Resize(n int) {
	*p = make(usersArray, n)
}
