// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
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

func (user *User) Markdown() md.MD {
	if user.mattermostUser != nil {
		return md.Markdownf("@%s", user.mattermostUser.Username)
	} else {
		return md.Markdownf("userID `%s`", user.MattermostUserID)
	}
}

func (user *User) MarkdownUnavailable(u *Unavailable) md.MD {
	return md.Markdownf("%s: %s", u.Reason, user.MarkdownInterval(u.Interval))
}

func (user *User) Time(t types.Time) types.Time {
	return t.In(user.location)
}

func (user *User) MarkdownInterval(i types.Interval) md.MD {
	return md.Markdownf("%s to %s",
		user.Time(i.Start), user.Time(i.Finish))
}

func (user *User) MarkdownWithSkills() md.MD {
	return md.Markdownf("%s %s", user.Markdown(), user.MarkdownSkills())
}

func (user *User) MarkdownSkills() md.MD {
	skillLevels := []string{}
	for _, s := range user.SkillLevels.IDs() {
		skillLevels = append(skillLevels, NewSkillLevel(s, Level(user.SkillLevels.Get(s))).String())
	}
	if len(skillLevels) == 0 {
		return "(none)"
	}
	ss := strings.Join(skillLevels, ", ")
	return md.Markdownf("(%s)", ss)
}

func (user User) MattermostUsername() string {
	if user.mattermostUser == nil {
		return string(user.MattermostUserID)
	}
	return user.mattermostUser.Username
}

func (user *User) AddUnavailable(uu ...*Unavailable) []*Unavailable {
	var added []*Unavailable
UNAVAILABLE:
	for _, u := range uu {
		for _, existing := range user.Calendar {
			if existing == u {
				continue UNAVAILABLE
			}
		}
		user.Calendar = append(user.Calendar, u)
		added = append(added, u)
	}
	unavailableBy(byStartDate).Sort(user.Calendar)
	return added
}

func (user *User) FindUnavailable(matchInterval types.Interval, matchRotationID, matchTaskID types.ID) []*Unavailable {
	var found []*Unavailable
	user.ScanUnavailable(
		matchInterval, matchRotationID, matchTaskID,
		func(event *Unavailable) {
			found = append(found, event)
		},
		nil,
	)
	return found
}

func (user *User) ClearUnavailable(matchInterval types.Interval, matchRotationID, matchTaskID types.ID) []*Unavailable {
	var cleared, kept []*Unavailable
	user.ScanUnavailable(
		matchInterval, matchRotationID, matchTaskID,
		func(event *Unavailable) {
			cleared = append(cleared, event)
		},
		func(event *Unavailable) {
			kept = append(kept, event)
		},
	)
	user.Calendar = kept
	return cleared
}

func (user *User) ScanUnavailable(matchInterval types.Interval, matchRotationID, matchTaskID types.ID,
	matchf, nonmatchf func(*Unavailable)) {
	for _, event := range user.Calendar {
		if !matchInterval.IsEmpty() && !event.Overlaps(matchInterval) {
			if nonmatchf != nil {
				nonmatchf(event)
			}
			continue
		}
		if matchTaskID != "" && event.TaskID != matchTaskID {
			if nonmatchf != nil {
				nonmatchf(event)
			}
			continue
		}
		if matchRotationID != "" {
			if (event.Reason != ReasonTask && event.Reason != ReasonGrace) ||
				(event.RotationID != matchRotationID) {
				if nonmatchf != nil {
					nonmatchf(event)
				}
			}
		}

		if matchf != nil {
			matchf(event)
		}
	}
}
