// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationSetAutopilot(parameters []string) (md.MD, error) {
	off := c.assureFS().Bool("off", false, "turn off")
	create := c.assureFS().Bool("create", false, "create shifts automatically")
	createPrior := c.assureFS().Duration("create-prior", 0, "create shifts this long before their scheduled start")
	schedule := c.assureFS().Bool("schedule", false, "create shifts automatically")
	schedulePrior := c.assureFS().Duration("schedule-prior", 0, "fill and schedule shifts this long before their scheduled start")
	startFinish := c.assureFS().Bool("start-finish", false, "start and finish scheduled tasks")
	remindStart := c.assureFS().Bool("remind-start", false, "remind shift users prior to start")
	remindStartPrior := c.assureFS().Duration("remind-start-prior", 0, "remind shift users this long before the shift's start")
	remindFinish := c.assureFS().Bool("remind-finish", false, "remind shift users prior to finish")
	remindFinishPrior := c.assureFS().Duration("remind-finish-prior", 0, "remind shift users this long before the shift's finish")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	var out md.Markdowner
	switch {
	case *off:
		out, err = c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.AutopilotSettings = sl.AutopilotSettings{}
			return nil
		})

	default:
		out, err = c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.AutopilotSettings.Create = *create
			r.AutopilotSettings.CreatePrior = *createPrior
			r.AutopilotSettings.Schedule = *schedule
			r.AutopilotSettings.SchedulePrior = *schedulePrior
			r.AutopilotSettings.StartFinish = *startFinish
			r.AutopilotSettings.RemindStart = *remindStart
			r.AutopilotSettings.RemindStartPrior = *remindStartPrior
			r.AutopilotSettings.RemindFinish = *remindFinish
			r.AutopilotSettings.RemindFinishPrior = *remindFinishPrior
			return nil
		})
	}

	return c.normalOut(out, err)
}
