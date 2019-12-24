// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func MarkdownRotation(rotation *Rotation) string {
	return fmt.Sprintf("%s", rotation.RotationID)
}

func MarkdownRotationWithDetails(rotation *Rotation) string {
	return fmt.Sprintf("**%s**: period %s, starting %s, %v users joined, needs %s",
		rotation.RotationID, rotation.Period, rotation.Start,
		len(rotation.MattermostUserIDs), MarkdownNeeds(rotation.Needs))
}

func MarkdownUserMapWithSkills(m UserMap) string {
	out := []string{}
	for _, user := range m {
		out = append(out, fmt.Sprintf("%s: %s", MarkdownUser(user), MarkdownUserSkills(user)))
	}
	return strings.Join(out, ", ")
}

func MarkdownShift(shiftNumber int, shift *Shift) string {
	return fmt.Sprintf("%v: %s to %s: %s",
		shiftNumber, shift.Start, shift.End, MarkdownUserMapWithSkills(shift.Users))
}
func MarkdownUserMap(m UserMap) string {
	out := []string{}
	for _, user := range m {
		out = append(out, MarkdownUser(user))
	}
	return strings.Join(out, ", ")
}

func MarkdownUserWithSkills(user *User) string {
	return fmt.Sprintf("%s: %s", MarkdownUser(user), MarkdownUserSkills(user))
}

func MarkdownUser(user *User) string {
	if user.MattermostUser != nil {
		return fmt.Sprintf("@%s", user.MattermostUser.Username)
	} else {
		return fmt.Sprintf("userID:`%s`", user.MattermostUserID)
	}
}

func MarkdownUserSkills(user *User) string {
	skills := []string{}
	for s, l := range user.SkillLevels {
		skills = append(skills, MarkdownSkillLevel(s, Level(l)))
	}

	if len(skills) == 0 {
		return "kook"
	}
	ss := strings.Join(skills, ", ")
	return fmt.Sprintf("(%s)", ss)
}

func MarkdownSkillLevel(skillName string, level Level) string {
	return fmt.Sprintf("%s%s", Level(level).String(), skillName)
}

func MarkdownNeed(need store.Need) string {
	if need.Max == 0 {
		return fmt.Sprintf("%v of %s", need.Min, MarkdownSkillLevel(need.Skill, Level(need.Level)))
	} else {
		return fmt.Sprintf("%v(%v) of %s", need.Min, need.Max, MarkdownSkillLevel(need.Skill, Level(need.Level)))
	}
}

func MarkdownNeeds(needs []store.Need) string {
	out := []string{}
	for _, need := range needs {
		out = append(out, MarkdownNeed(need))
	}
	return strings.Join(out, ", ")
}

func MarkdownNeedsList(needs map[string]store.Need, indent string) string {
	out := ""
	for _, need := range needs {
		out += indent + "- " + MarkdownNeed(need) + "\n"
	}
	return out
}
