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
	SkillLevels      *types.IntSet  `json:",omitempty"` // skill (id) -> level
	LastServed       *types.IntSet  `json:",omitempty"` // Last time completed a task, rotationID -> Unix time.
	Calendar         []*Unavailable `json:",omitempty"` // Sorted by start date of the events.

	// private fields
	loaded         bool
	mattermostUser *model.User
	location       *time.Location
}

func NewUser(mattermostUserID types.ID) *User {
	return &User{
		MattermostUserID: mattermostUserID,
		Calendar:         []*Unavailable{},
		SkillLevels:      types.NewIntSet(),
		LastServed:       types.NewIntSet(),
	}
}

func (user *User) GetID() types.ID {
	return user.MattermostUserID
}

func (user *User) WithLastServed(rotationID types.ID, finishTime types.Time) *User {
	user.LastServed.Set(rotationID, finishTime.Unix())
	return user
}

func (user *User) WithSkills(skillLevels *types.IntSet) *User {
	user.SkillLevels = skillLevels
	return user
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

func (user *User) FindUnavailable(interval types.Interval, applicableToRotationID types.ID) []*Unavailable {
	var found []*Unavailable
	for _, unavailable := range user.Calendar {
		s, f := unavailable.Start, unavailable.Finish
		if s.Before(interval.Start.Time) {
			s = interval.Start
		}
		if f.After(interval.Finish.Time) {
			f = interval.Finish
		}
		if !s.Before(f.Time) {
			continue
		}

		// Overlap, only consider events applicable to the rotation
		if applicableToRotationID == "" ||
			(unavailable.Reason != ReasonTask && unavailable.Reason != ReasonGrace) ||
			unavailable.RotationID != applicableToRotationID {
			continue
		}

		found = append(found, unavailable)
	}
	return found
}

func (user *User) ClearUnavailable(interval types.Interval, applicableToRotationID types.ID) []*Unavailable {
	var cleared, updated []*Unavailable
	for _, unavailable := range user.Calendar {
		s, f := unavailable.Start, unavailable.Finish
		if s.Before(interval.Start.Time) {
			s = interval.Start
		}
		if f.After(interval.Finish.Time) {
			f = interval.Finish
		}

		if s.Before(f.Time) {
			// Overlap, only consider events applicable to the rotation
			if applicableToRotationID == "" {
				cleared = append(cleared, unavailable)
				continue
			}
			if (unavailable.Reason == ReasonTask || unavailable.Reason == ReasonGrace) &&
				unavailable.RotationID == applicableToRotationID {
				cleared = append(cleared, unavailable)
				continue
			}
		}

		updated = append(updated, unavailable)
	}
	user.Calendar = updated
	return cleared
}

func (user *User) GetLastServed(r *Rotation) int64 {
	last := user.LastServed.Get(r.RotationID)
	if last <= 0 {
		return r.Beginning.Unix()
	}
	return last
}
