// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"fmt"
	"strconv"
	"strings"
)

type Needs []*Need

type Need struct {
	Min   int
	Max   int
	Skill string
	Level int
}

func NewNeed(skill string, level int, min int) *Need {
	return &Need{
		Min:   min,
		Max:   -1,
		Skill: skill,
		Level: level,
	}
}

func (need Need) WithMax(max int) *Need {
	need.Max = max
	return &need
}

func (need Need) String() string {
	return fmt.Sprintf("%s-%v-%v(%v)", need.Skill, need.Level, need.Min, need.Max)
}

func (need Need) SkillLevel() string {
	return need.Skill + "-" + strconv.Itoa(need.Level)
}

func (need *Need) Markdown() string {
	if need.Max == -1 {
		return fmt.Sprintf("**%v** %s", need.Min, need.SkillLevel())
	} else {
		return fmt.Sprintf("**%v(%v)** %s", need.Min, need.Max, need.SkillLevel())
	}
}

func (needs Needs) Clone() Needs {
	var newNeeds Needs
	for _, need := range needs {
		newNeed := *need
		newNeeds = append(newNeeds, &newNeed)
	}
	return newNeeds
}

func (needs Needs) Markdown() string {
	out := []string{}
	for _, need := range needs {
		out = append(out, need.Markdown())
	}
	return strings.Join(out, ", ")
}

func (needs Needs) MarkdownBullets() string {
	out := ""
	for _, need := range needs {
		out += "- " + need.Markdown() + "\n"
	}
	return out
}
