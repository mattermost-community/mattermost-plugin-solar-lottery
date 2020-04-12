// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationAutopilot(parameters []string) (md.MD, error) {
	off := c.withFlagOff()
	create, createPrior := c.withFlagCreatePrior()
	schedule, schedulePrior := c.withFlagSchedulePrior()
	startFinish := c.withFlagStartFinish()
	remindStart, remindStartPrior := c.withFlagRemindStartPrior()
	remindFinish, remindFinishPrior := c.withFlagRemindStartPrior()
	run := c.withFlagRun()
	now, err := c.withFlagNow()
	if err != nil {
		return c.flagUsage(), err
	}
	err = c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	var out md.Markdowner
	switch {
	case *run:
		out, err = c.SL.RunAutopilot(&sl.InRunAutopilot{
			RotationID: rotationID,
			Time:       *now,
		})

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
