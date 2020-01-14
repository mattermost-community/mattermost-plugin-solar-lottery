// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

func (api *api) AutopilotRotation(rotation *Rotation, now time.Time) error {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.AutopilotRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
		"Time":           now,
	})

	if !rotation.Autopilot.On {
		return nil
	}
	currentShiftNumber, err := rotation.ShiftNumberForTime(now)
	if err != nil {
		return err
	}

	api.Debugf("running autopilot for shiftNumber: %v", currentShiftNumber)
	status := func(err error, message ...string) string {
		switch {
		case err == nil && len(message) > 0 && message[0] != "":
			return "ok: " + message[0]
		case err == nil && (len(message) == 0 || message[0] == ""):
			return "ok"
		case err != nil:
			return "**" + err.Error() + "**"
		}
		return "unknown"
	}

	finishedShiftNumber, finishedShift, err :=
		api.autopilotFinishShift(rotation, now, currentShiftNumber)
	finishedStatus := status(err, fmt.Sprintf("finished %s.\n%s",
		api.MarkdownShift(rotation, finishedShiftNumber),
		api.MarkdownIndent(api.MarkdownShiftBullets(rotation, finishedShiftNumber, finishedShift), "  ")))

	filledShiftNumbers, filledShifts, filledAdded, err :=
		api.autopilotFill(rotation, now, currentShiftNumber, logger)
	fillStatus := fmt.Sprintf("ok: **processed %v shifts**", len(filledShifts))
	if err != nil {
		if len(filledShifts) > 0 {
			fillStatus = fmt.Sprintf("error: processed %v shifts, then **failed: %v**.", len(filledShifts), err)
		} else {
			fillStatus = err.Error()
		}
	}
	if len(filledShifts) > 0 {
		fillStatus += "\n"
	}
	for i, shift := range filledShifts {
		if len(filledAdded[i]) > 0 {
			fillStatus += fmt.Sprintf(
				"  - %s: **added users** %s.\n%s",
				api.MarkdownShift(rotation, filledShiftNumbers[i]),
				api.MarkdownUsers(filledAdded[i]),
				api.MarkdownIndent(api.MarkdownShiftBullets(rotation, filledShiftNumbers[i], shift), "    "))
		} else {
			fillStatus += fmt.Sprintf(
				"  - %s: no change.\n",
				api.MarkdownShift(rotation, filledShiftNumbers[i]))
		}
	}

	startedShiftNumber, startedShift, err := api.autopilotStartShift(rotation, now, currentShiftNumber)
	startedStatus := status(err, fmt.Sprintf("started %s.\n%s",
		api.MarkdownShift(rotation, startedShiftNumber),
		api.MarkdownIndent(api.MarkdownShiftBullets(rotation, startedShiftNumber, startedShift), "  ")))

	currentNotified, err := api.autopilotNotifyCurrent(rotation, now, currentShiftNumber)
	currentNotifiedStatus := status(err, fmt.Sprintf("notified %s", api.MarkdownUsers(currentNotified)))

	nextNotified, err := api.autopilotNotifyNext(rotation, now, currentShiftNumber)
	nextNotifiedStatus := status(err, fmt.Sprintf("notified %s", api.MarkdownUsers(nextNotified)))

	logger.Infof("%s ran autopilot on %s for %v. Status:\n"+
		"- finish previous shift: %s\n"+
		"- fill shift(s): %s\n"+
		"- start next shift: %s\n"+
		"- notify current shift's users: %s\n"+
		"- notify next shift's users: %s\n",
		api.MarkdownUser(api.actingUser), api.MarkdownRotation(rotation), now,
		finishedStatus,
		fillStatus,
		startedStatus,
		currentNotifiedStatus,
		nextNotifiedStatus,
	)

	return nil
}
