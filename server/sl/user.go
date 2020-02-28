// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type User struct {
	PluginVersion    string `json:",omitempty"`
	MattermostUserID types.ID
	SkillLevels      *types.IntIndex `json:",omitempty"` // skill (id) -> level
	LastServed       *types.IntIndex `json:",omitempty"` // Last time completed a task, rotationID -> Unix time.
	Calendar         []*Unavailable  `json:",omitempty"` // Sorted by start date of the events.

	// private fields
	loaded         bool
	mattermostUser *model.User
	location       *time.Location
}

func NewUser(mattermostUserID types.ID) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		SkillLevels:      types.NewIntIndex(),
		LastServed:       types.NewIntIndex(),
		Calendar:         []*Unavailable{},
	}
}

func (u *User) CloneUser(deep bool) *User {
	return u.Clone(deep).(*User)
}

func (u *User) Clone(deep bool) types.Cloneable {
	clone := *u
	clone.SkillLevels = u.SkillLevels.Clone(deep).(*types.IntIndex)
	clone.LastServed = u.LastServed.Clone(deep).(*types.IntIndex)
	clone.Calendar = append([]*Unavailable{}, u.Calendar...)
	return &clone
}

func (user *User) WithLastServed(rotationID types.ID, finishTime types.Time) *User {
	newUser := user.CloneUser(false)
	newUser.LastServed.Set(rotationID, finishTime.Unix())
	return newUser
}

func (user *User) WithSkills(skillsLevels *types.IntIndex) *User {
	newUser := user.Clone(false).(*User)
	newUser.SkillLevels = types.NewIntIndex()
	for _, s := range skillsLevels.IDs() {
		newUser.SkillLevels.Set(s, skillsLevels.Get(s))
	}
	return newUser
}

func (user *User) String() string {
	if user.mattermostUser != nil {
		return fmt.Sprintf("@%s", user.mattermostUser.Username)
	} else {
		return fmt.Sprintf("%q", user.MattermostUserID)
	}
}

func (user *User) Markdown() string {
	if user.mattermostUser != nil {
		return fmt.Sprintf("@%s", user.mattermostUser.Username)
	} else {
		return fmt.Sprintf("userID `%s`", user.MattermostUserID)
	}
}

func (user *User) MarkdownUnavailable(u *Unavailable) string {
	return fmt.Sprintf("%s: %s", u.Reason, user.MarkdownInterval(u.Interval))
}

func (user *User) Time(t types.Time) types.Time {
	return t.In(user.location)
}

func (user *User) MarkdownInterval(i types.Interval) string {
	return fmt.Sprintf("%s to %s",
		user.Time(i.Start), user.Time(i.Finish))
}

func (user *User) MarkdownWithSkills() string {
	return fmt.Sprintf("%s %s", user.Markdown(), user.MarkdownSkills())
}

func (user *User) MarkdownSkills() string {
	skillLevels := []string{}
	for _, s := range user.SkillLevels.IDs() {
		skillLevels = append(skillLevels, NewSkillLevel(s, Level(user.SkillLevels.Get(s))).String())
	}
	if len(skillLevels) == 0 {
		return "(none)"
	}
	ss := strings.Join(skillLevels, ", ")
	return fmt.Sprintf("(%s)", ss)
}

func (user User) MattermostUsername() string {
	if user.mattermostUser == nil {
		return string(user.MattermostUserID)
	}
	return user.mattermostUser.Username
}

func (user *User) AddUnavailable(uu ...*Unavailable) {
UNAVAILABLE:
	for _, u := range uu {
		for _, existing := range user.Calendar {
			if existing == u {
				continue UNAVAILABLE
			}
		}
		user.Calendar = append(user.Calendar, u)
		unavailableBy(byStartDate).Sort(user.Calendar)
	}
}

func (user *User) findUnavailable(interval types.Interval, remove bool) []*Unavailable {
	var found, updated []*Unavailable
	for _, u := range user.Calendar {
		s, f := u.Start, u.Finish
		if s.Before(interval.Start.Time) {
			s = interval.Start
		}
		if f.After(interval.Finish.Time) {
			f = interval.Finish
		}

		if s.Before(f.Time) {
			// Overlap
			found = append(found, u)
			if remove {
				continue
			}
		}

		updated = append(updated, u)
	}
	user.Calendar = updated
	return found
}

func (user *User) IsQualified(skillLevel SkillLevel) bool {
	level := user.SkillLevels.Get(skillLevel.Skill)
	return level >= int64(skillLevel.Level)
}
