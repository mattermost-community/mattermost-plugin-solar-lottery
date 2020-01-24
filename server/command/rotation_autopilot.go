// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"time"

	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func withRotationAutopilotFlags(fs *pflag.FlagSet, off *bool, autostart *bool, autofill *bool, autofillPriorDays *int, notifyPriorDays *int, debugRunTime *string) {
	fs.StringVar(debugRunTime, flagDebugRun, "", "run autopilot mocking the specified time")
	fs.IntVar(notifyPriorDays, flagNotifyDays, 7, "notify shift users this many days prior to transition date")
	fs.IntVar(autofillPriorDays, flagFillDays, 30, "autofill shifts this many days prior to start date")
	fs.BoolVar(autostart, flagStart, true, "start and finish shifts automatically")
	fs.BoolVar(autofill, flagFill, true, "start and finish shifts automatically")
	fs.BoolVar(off, flagOff, false, "turn autopilot off")
}

func (c *Command) autopilotRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	var autostart, autofill, off bool
	var autofillPriorDays, notifyPriorDays int
	var debugRunTime string
	fs := newRotationFlagSet(&rotationID, &rotationName)
	withRotationAutopilotFlags(fs, &off, &autostart, &autofill, &autofillPriorDays, &notifyPriorDays, &debugRunTime)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.SL.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	if debugRunTime != "" {
		var now time.Time
		now, err = time.Parse(sl.DateFormat, debugRunTime)
		if err != nil {
			return c.flagUsage(fs), err
		}

		err = c.SL.AutopilotRotation(rotation, now)
		if err != nil {
			return "", err
		}
		return "Ran autopilot for " + now.String(), nil
	}

	var updatef func(rotation *sl.Rotation) error
	if off {
		updatef = func(rotation *sl.Rotation) error {
			rotation.Autopilot = store.RotationAutopilot{}
			return nil
		}
	} else {
		updatef = func(rotation *sl.Rotation) error {
			rotation.Autopilot = store.RotationAutopilot{
				On:          true,
				StartFinish: autostart,
				Fill:        autofill,
				FillPrior:   time.Duration(autofillPriorDays) * sl.DayDuration,
				Notify:      notifyPriorDays != 0,
				NotifyPrior: time.Duration(notifyPriorDays) * sl.DayDuration,
			}
			return nil
		}
	}

	err = c.SL.UpdateRotation(rotation, updatef)
	if err != nil {
		return "", err
	}

	return "Updated rotation autopilot:\n" + rotation.MarkdownBullets(), nil
}
