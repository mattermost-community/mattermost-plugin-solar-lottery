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
	MattermostUserID string
	SkillLevels      types.IntMap   `json:",omitempty"`
	LastServed       types.IntMap   `json:",omitempty"` // Last time completed a task, per rotation ID; Unix time.
	Calendar         []*Unavailable `json:",omitempty"` // Sorted by start date of the events.

	// private fields
	loaded         bool
	mattermostUser *model.User
	location       *time.Location
}

func NewUser(mattermostUserID string) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		SkillLevels:      types.IntMap{},
		LastServed:       types.IntMap{},
		Calendar:         []*Unavailable{},
	}
}

func (u *User) Clone() *User {
	clone := *u
	clone.SkillLevels = u.SkillLevels.Clone()
	clone.LastServed = u.LastServed.Clone()
	clone.Calendar = append([]*Unavailable{}, u.Calendar...)
	return &clone
}

func (user *User) WithLastServed(rotationID string, finishTime types.Time) *User {
	newUser := user.Clone()
	newUser.LastServed[rotationID] = finishTime.Unix()
	return newUser
}

func (user *User) WithSkills(skillsLevels types.IntMap) *User {
	newUser := user.Clone()
	if newUser.SkillLevels != nil {
		newUser.SkillLevels = types.IntMap{}
	}
	for s, l := range skillsLevels {
		newUser.SkillLevels[s] = l
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
	skills := []string{}
	for s, l := range user.SkillLevels {
		skills = append(skills, MarkdownSkillLevel(s, Level(l)))
	}

	if len(skills) == 0 {
		return "(kook)"
	}
	ss := strings.Join(skills, ", ")
	return fmt.Sprintf("(%s)", ss)
}

func (user User) MattermostUsername() string {
	if user.mattermostUser == nil {
		return user.MattermostUserID
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

func (user *User) IsQualified(need *Need) bool {
	skillLevel, _ := user.SkillLevels[need.Skill]
	return skillLevel >= int64(need.Level)
}
