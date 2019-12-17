// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

// syntax:
// - `...skill list [<skill-name>]`
// 		List skills, and the number of skilled users by level.
// - `...skill add skill-name [--users [@user1[ ,@user2...]] [--level]`
//		Adds a skill level to user(s). If the skill is not known, adds the
//		skill as a known skill.
// - `...skill delete skill-name [--users @user1[ ,@user2...]]`
//		Deletes a skill from the specified users. If --users is provided,
// 		deletes the skill altogether.
func (c *Command) skill(parameters ...string) (string, error) {
	fs := flag.NewFlagSet("skill", flag.ContinueOnError)
	users := ""
	level := ""
	fs.StringVar(&users, "users", "", "a list of users")
	fs.StringVar(&level, "level", api.LevelBeginner, "beginner, intermediate, advanced, or expert")
	err := fs.Parse(parameters)
	if err != nil {
		return "", err
	}

	args := fs.Args()
	if len(args) < 1 {
		return "", errors.New("must provide an action ")
	}
	cmd := fs.Arg(0)
	skillName := fs.Arg(1)

	if len(users) > 0 {
		return c.skillUsers(cmd, skillName, level, users)
	} else {
		return c.skillNoUsers(cmd, skillName)
	}
}

func (c *Command) skillNoUsers(cmd, skillName string) (string, error) {
	var err error
	out := ""
	switch cmd {
	case "add":
		err = c.API.AddSkill(skillName)
		out = fmt.Sprintf("Added **%s** to known skills: ", skillName)

	case "delete":
		err = c.API.DeleteSkill(skillName)
		out = fmt.Sprintf("Removed **%s** from known skills: ", skillName)

	case "list":
		var skills store.IDMap
		skills, err = c.API.ListSkills()
		out = "Known skills: " + utils.JSONBlock(skills)

	default:
		return "", errors.Errorf("%q is not valid, must be add, delete, list", cmd)
	}
	if err != nil {
		return "", err
	}
	return out, nil
}

func (c *Command) skillUsers(cmd, skillName, level, users string) (string, error) {
	var err error
	switch cmd {
	case "add":
		l := 1
		if level != "" {
			l, err = api.ParseLevel(level)
			if err != nil {
				return "", err
			}
		}
		err = c.API.AddSkillToUsers(users, skillName, l)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Added/updated %s(%v) to %s", skillName, level, users), nil

	case "delete":
	case "list":
	}
	return "", nil
}
