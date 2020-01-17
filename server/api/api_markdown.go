// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

type Markdowner interface {
	MarkdownEvent(event Event) string
	MarkdownIndent(in, prefix string) string
	MarkdownNeed(need store.Need) string
	MarkdownNeeds(needs []store.Need) string
	MarkdownNeedsBullets(needs map[string]store.Need, indent string) string
	MarkdownRotation(rotation *Rotation) string
	MarkdownRotationBullets(*Rotation) string
	MarkdownShift(rotation *Rotation, shiftNumber int) string
	MarkdownShiftBullets(*Rotation, int, *Shift) string
	MarkdownSkillLevel(skillName string, level Level) string
	MarkdownUser(*User) string
	MarkdownUsers(UserMap) string
	MarkdownUserSkills(user *User) string
	MarkdownUsersWithSkills(UserMap) string
	MarkdownUserWithSkills(user *User) string
}

func (api *api) MarkdownRotation(rotation *Rotation) string {
	return markdownRotation(rotation)
}

func markdownRotation(rotation *Rotation) string {
	return fmt.Sprintf("%s", rotation.Name)
}

func (api *api) MarkdownEvent(event Event) string {
	return fmt.Sprintf("%s: %s to %s",
		event.Type, event.Start, event.End)
}

func (api *api) MarkdownShift(rotation *Rotation, shiftNumber int) string {
	return fmt.Sprintf("%s#%v", rotation.Name, shiftNumber)
}

func (api *api) MarkdownUser(user *User) string {
	api.ExpandUser(user)
	if user.MattermostUser != nil {
		return fmt.Sprintf("@%s", user.MattermostUser.Username)
	} else {
		return fmt.Sprintf("userID `%s`", user.MattermostUserID)
	}
}

func (api *api) MarkdownRotationBullets(rotation *Rotation) string {
	api.ExpandRotation(rotation)

	out := fmt.Sprintf("- **%s**\n", rotation.Name)
	out += fmt.Sprintf("  - ID: `%s`.\n", rotation.RotationID)
	out += fmt.Sprintf("  - Starting: **%s**.\n", rotation.Start)
	out += fmt.Sprintf("  - Period: **%s**.\n", rotation.Period)
	out += fmt.Sprintf("  - Size: **%v** people.\n", rotation.Size)
	out += fmt.Sprintf("  - Needs (%v): %s.\n", len(rotation.Needs), api.MarkdownNeeds(rotation.Needs))
	out += fmt.Sprintf("  - Grace: **%v** shifts.\n", rotation.Grace)
	out += fmt.Sprintf("  - Users (%v): %s.\n", len(rotation.MattermostUserIDs), api.MarkdownUsersWithSkills(rotation.Users))

	if rotation.Autopilot.On {
		out += fmt.Sprintf("  - Autopilot: **on**\n")
		out += fmt.Sprintf("    - Auto-start: **%v**\n", rotation.Autopilot.StartFinish)
		out += fmt.Sprintf("    - Auto-fill: **%v**, %v days prior to start\n", rotation.Autopilot.Fill, rotation.Autopilot.FillPrior)
		out += fmt.Sprintf("    - Notify users in advance: **%v**, %v days prior to transition\n", rotation.Autopilot.Notify, rotation.Autopilot.NotifyPrior)
	} else {
		out += fmt.Sprintf("  - Autopilot: **off**\n")
	}

	return out
}

func (api *api) MarkdownUsersWithSkills(m UserMap) string {
	out := []string{}
	for _, user := range m {
		out = append(out, fmt.Sprintf("%s %s", api.MarkdownUser(user), api.MarkdownUserSkills(user)))
	}
	return strings.Join(out, ", ")
}

func (api *api) MarkdownShiftBullets(rotation *Rotation, shiftNumber int, shift *Shift) string {
	if shift == nil {
		return "n/a"
	}
	api.ExpandRotation(rotation)

	out := fmt.Sprintf("- **%s**: %s to %s\n", api.MarkdownShift(rotation, shiftNumber), shift.Start, shift.End)
	out += fmt.Sprintf("  - Status: **%s**\n", shift.Status)
	out += fmt.Sprintf("  - Users: **%v**\n", len(shift.MattermostUserIDs))
	for _, user := range rotation.ShiftUsers(shift) {
		out += fmt.Sprintf("    - %s\n", api.MarkdownUserWithSkills(user))
	}
	return out
}

func (api *api) MarkdownUsers(m UserMap) string {
	out := []string{}
	for _, user := range m {
		out = append(out, api.MarkdownUser(user))
	}
	return strings.Join(out, ", ")
}

func (api *api) MarkdownUserWithSkills(user *User) string {
	return fmt.Sprintf("%s %s", api.MarkdownUser(user), api.MarkdownUserSkills(user))
}

func (api *api) MarkdownUserSkills(user *User) string {
	skills := []string{}
	for s, l := range user.SkillLevels {
		skills = append(skills, api.MarkdownSkillLevel(s, Level(l)))
	}

	if len(skills) == 0 {
		return "(kook)"
	}
	ss := strings.Join(skills, ", ")
	return fmt.Sprintf("(%s)", ss)
}

func (api *api) MarkdownSkillLevel(skillName string, level Level) string {
	return markdownSkillLevel(skillName, level)
}

func markdownSkillLevel(skillName string, level Level) string {
	return fmt.Sprintf("%s%s", Level(level).String(), skillName)
}

func (api *api) MarkdownNeed(need store.Need) string {
	return markdownNeed(need)
}

func markdownNeed(need store.Need) string {
	if need.Max == -1 {
		return fmt.Sprintf("**%v** %s", need.Min, markdownSkillLevel(need.Skill, Level(need.Level)))
	} else {
		return fmt.Sprintf("**%v(%v)** %s", need.Min, need.Max, markdownSkillLevel(need.Skill, Level(need.Level)))
	}
}

func (api *api) MarkdownNeeds(needs []store.Need) string {
	out := []string{}
	for _, need := range needs {
		out = append(out, api.MarkdownNeed(need))
	}
	return strings.Join(out, ", ")
}

func (api *api) MarkdownNeedsBullets(needs map[string]store.Need, indent string) string {
	out := ""
	for _, need := range needs {
		out += indent + "- " + api.MarkdownNeed(need) + "\n"
	}
	return out
}

func (api *api) MarkdownIndent(in, prefix string) string {
	lines := strings.Split(in, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
