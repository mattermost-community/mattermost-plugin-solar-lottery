// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationSetTask(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	dur := c.assureFS().Duration("duration", 0, "duration")
	grace := c.assureFS().Duration("grace", 0, "grace period after finishing a task")
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
			if *dur != 0 {
				r.TaskSettings.Duration = *dur
			}
			if *grace != 0 {
				r.TaskSettings.Grace = *grace
			}
			return nil
		}))
}
