// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) skill(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandAdd:    c.addSkill,
		commandDelete: c.deleteSkill,
		commandList:   c.listSkills,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) addSkill(parameters []string) (string, error) {
	var skillName string
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	err = c.SL.AddKnownSkill(skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Added **%s** to known skills.", skillName), nil
}

func (c *Command) deleteSkill(parameters []string) (string, error) {
	var skillName string
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withSkillFlags(fs, &skillName, nil)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	err = c.SL.DeleteSkill(skillName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted **%s** from known skills. User profiles are not changed.", skillName), nil
}

func (c *Command) listSkills(parameters []string) (string, error) {
	skills, err := c.SL.ListSkills()
	if err != nil {
		return "", err
	}
	return "Known skills: " + utils.JSONBlock(skills), nil
}

func withSkillFlags(fs *pflag.FlagSet, skillName *string, level *sl.Level) {
	fs.StringVarP(skillName, flagSkill, flagPSkill, "", "Skill name")
	if level != nil {
		fs.VarP(level, flagLevel, flagPLevel, "Skill level")
	}
}
