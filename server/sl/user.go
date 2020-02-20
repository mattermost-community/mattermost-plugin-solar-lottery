// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type User struct {
	PluginVersion    string `json:",omitempty"`
	MattermostUserID string
	SkillLevels      IntMap         `json:",omitempty"`
	LastServed       IntMap         `json:",omitempty"` // Last time completed a task, per rotation ID; Unix time.
	Calendar         []*Unavailable `json:",omitempty"` // Sorted by start date of the events.

	// private fields
	loaded         bool
	mattermostUser *model.User
}

func NewUser(mattermostUserID string) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		SkillLevels:      IntMap{},
		LastServed:       IntMap{},
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

func (user *User) WithSkills(skillsLevels IntMap) *User {
	newUser := user.Clone()
	if newUser.SkillLevels != nil {
		newUser.SkillLevels = IntMap{}
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

// func (user *User) OverlapEvents(intervalStart, intervalEnd time.Time, remove bool) ([]store.Event, error) {
// 	var found, updated []store.Event
// 	for _, event := range user.Events {
// 		s, e, err := ParseDatePair(event.Start, event.End)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Find the overlap
// 		if s.Before(intervalStart) {
// 			s = intervalStart
// 		}
// 		if e.After(intervalEnd) {
// 			e = intervalEnd
// 		}

// 		if s.Before(e) {
// 			// Overlap
// 			found = append(found, event)
// 			if remove {
// 				continue
// 			}
// 		}

// 		updated = append(updated, event)
// 	}
// 	user.Events = updated
// 	return found, nil
// }

func (user *User) IsQualified(need *Need) bool {
	skillLevel, _ := user.SkillLevels[need.Skill]
	return skillLevel >= int64(need.Level)
}
