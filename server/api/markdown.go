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
	out := fmt.Sprintf("###### %s\n", rotation.Name)
	out += fmt.Sprintf("- ID: `%s`\n", rotation.RotationID)
	out += fmt.Sprintf("- Starting: `%s`\n", rotation.Start)
	out += fmt.Sprintf("- Period: `%s`\n", rotation.Period)
	out += fmt.Sprintf("- Needs: %s\n", MarkdownNeeds(rotation.Needs))
	out += fmt.Sprintf("- Grace: `%v`\n", rotation.Grace)
	out += fmt.Sprintf("- Users (%v): %s\n", len(rotation.MattermostUserIDs), MarkdownUserMapWithSkills(rotation.Users))

	if rotation.Autopilot.On {
		out += fmt.Sprintf("- Autopilot: `on`\n")
		out += fmt.Sprintf("  - Auto-start: `%v`\n", rotation.Autopilot.StartFinish)
		out += fmt.Sprintf("  - Auto-fill: `%v`, %v prior to start\n", rotation.Autopilot.Fill, rotation.Autopilot.FillPrior)
		out += fmt.Sprintf("  - Notify users: `%v`, %v prior to transition\n", rotation.Autopilot.Notify, rotation.Autopilot.NotifyPrior)
	} else {
		out += fmt.Sprintf("- Autopilot: `off`\n")
	}

	return out
}

func MarkdownUserMapWithSkills(m UserMap) string {
	out := []string{}
	for _, user := range m {
		out = append(out, fmt.Sprintf("%s: %s", MarkdownUser(user), MarkdownUserSkills(user)))
	}
	return strings.Join(out, ", ")
}

func MarkdownShift(rotation *Rotation, shiftNumber int, shift *Shift) string {
	return fmt.Sprintf("rotation %s shift #%v (%s to %s), status:%s, users: %s",
		rotation.Name, shiftNumber, shift.Start, shift.End, shift.Status, MarkdownUserMapWithSkills(shift.Users))
}

func MarkdownEvent(event store.Event) string {
	return fmt.Sprintf("%s: %s to %s",
		event.Type, event.Start, event.End)
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
