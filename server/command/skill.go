// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) addSkill(parameters []string) (string, error) {
	var skillName string
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
	}
	err = c.API.AddSkill(skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Added **%s** to known skills.", skillName), nil
}

func (c *Command) deleteSkill(parameters []string) (string, error) {
	var skillName string
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
	}
	err = c.API.DeleteSkill(skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted **%s** from known skills. User profiles are not changed.", skillName), nil
}

func (c *Command) listSkills(parameters []string) (string, error) {
	skills, err := c.API.ListSkills()
	if err != nil {
		return "", err
	}
	return "Known skills: " + utils.JSONBlock(skills), nil
}

func (c *Command) qualifyUsers(parameters []string) (string, error) {
	var usernames, skillName string
	var level api.Level
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	withSkillFlags(fs, &skillName, &level)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
	}

	err = c.API.AddSkillToUsers(usernames, skillName, level)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Qualified %s as %s", usernames, api.MarkdownSkillLevel(skillName, level)), nil
}

func (c *Command) disqualifyUsers(parameters []string) (string, error) {
	var usernames, skillName string
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to disqualify from skill")
	err := fs.Parse(parameters)
	if err != nil {
		return c.subUsage(fs), err
	}

	err = c.API.DeleteSkillFromUsers(usernames, skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Disqualified %s from %s", usernames, skillName), nil
}

func withSkillFlags(fs *pflag.FlagSet, skillName *string, level *api.Level) {
	fs.StringVarP(skillName, flagSkill, flagPSkill, "", "Skill name")
	if level != nil {
		fs.VarP(level, flagLevel, flagPLevel, "Skill level")
	}
}
