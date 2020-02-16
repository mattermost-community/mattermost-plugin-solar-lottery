// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

type Need struct {
	Count int
	Skill string
	Level int
}

type Needs []*Need

func NewNeed(skill string, level int, count int) *Need {
	return &Need{
		Count: count,
		Skill: skill,
		Level: level,
	}
}
