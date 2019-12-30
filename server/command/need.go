// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (c *Command) need(parameters []string) (string, error) {
	var rotationID, rotationName, skill string
	var level api.Level
	var deleteNeed bool
	var min, max int
	fs := flag.NewFlagSet("updateRotationNeed", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	withRotationNeedFlags(fs, &skill, &level, &min, &max, &deleteNeed)
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("show rotation", fs), err
	}
	if level == 0 || skill == "" {
		return subusage("rotation need", fs),
			errors.Errorf("requires `%s` and `%s` to be specified.", flagSkill, flagLevel)
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	// default to delete need
	updatef := func(rotation *api.Rotation) error {
		return rotation.DeleteNeed(skill, level)
	}
	if !deleteNeed {
		if min == 0 {
			return subusage("rotation need", fs),
				errors.Errorf("requires `%s` to be specified.", flagMin)
		}
		updatef = func(rotation *api.Rotation) error {
			rotation.ChangeNeed(skill, level, store.Need{
				Skill: skill,
				Level: int(level),
				Min:   min,
				Max:   max,
			})
			return nil
		}
	}

	err = c.API.UpdateRotation(rotation, updatef)
	if err != nil {
		return "", err
	}

	return "Updated rotation " + api.MarkdownRotationWithDetails(rotation), nil
}

func withRotationNeedFlags(fs *pflag.FlagSet, skill *string, level *api.Level, min, max *int, deleteNeed *bool) {
	fs.StringVar(skill, flagSkill, "", "if used with --need, indicates the needed skill")
	fs.Var(level, flagLevel, "if used with --need, indicates the needed skill level")
	fs.IntVar(min, flagMin, 0, "if used with --need, indicates the minimum needed headcount")
	fs.IntVar(max, flagMax, 0, "if used with --need, indicates the maximum needed headcount")
	fs.BoolVar(deleteNeed, flagDeleteNeed, false, "remove a need from rotation")
}