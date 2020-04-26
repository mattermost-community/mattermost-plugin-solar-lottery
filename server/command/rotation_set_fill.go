// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationSetFill(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	seed := c.assureFS().Int64("seed", intNoValue, "seed to use")
	fuzz := c.assureFS().Int64("fuzz", intNoValue, `increase fill randomness`)
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			if *seed != intNoValue {
				r.FillSettings.Seed = *seed
			}
			if *fuzz != intNoValue {
				r.FillSettings.Fuzz = *fuzz
			}
			return nil
		}))
}

func (c *Command) rotationSetRequire(parameters []string) (md.MD, error) {
	return c.rotationSetNeed(true, parameters)
}

func (c *Command) rotationSetLimit(parameters []string) (md.MD, error) {
	return c.rotationSetNeed(false, parameters)
}

func (c *Command) rotationSetNeed(require bool, parameters []string) (md.MD, error) {
	c.withFlagRotation()
	var skillLevel sl.SkillLevel
	c.assureFS().VarP(&skillLevel, "skill", "s", "skill, with optional level (1-4) as in `--skill=web-3`.")
	count := c.assureFS().Int("count", 1, "number of users")
	clear := c.assureFS().Bool("clear", false, "remove the skill from the list")
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			needsToUpdate := r.TaskSettings.Limit
			if require {
				needsToUpdate = r.TaskSettings.Require
			}
			if *clear {
				needsToUpdate.Delete(skillLevel.AsID())
			} else {
				needsToUpdate.SetCountForSkillLevel(skillLevel, int64(*count))
			}
			return nil
		}))
}
