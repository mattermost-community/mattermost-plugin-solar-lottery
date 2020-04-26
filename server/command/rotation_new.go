// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/queue"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) rotationNew(parameters []string) (md.MD, error) {
	fillType := c.assureFS().String("fill-type", "", fmt.Sprintf("fill type: %s or %s", solarlottery.Type, queue.Type))
	taskType := c.assureFS().String("task-type", "", fmt.Sprintf("task type: %s or %s", sl.TaskTypeShift, sl.TaskTypeTicket))
	beginning, err := c.withTimeFlag("beginning", "beginning of time for shifts")
	if err != nil {
		return c.flagUsage(), err
	}
	period := types.Period{}
	c.assureFS().Var(&period, "period", "recurrence period")
	seed := c.assureFS().Int64("seed", intNoValue, "seed to use")
	fuzz := c.assureFS().Int64("fuzz", intNoValue, `increase fill randomness`)

	err = c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	if c.fs.Arg(0) == "" {
		return c.flagUsage(), errors.Errorf("must specify rotation name")
	}

	switch types.ID(*fillType) {
	case solarlottery.Type, queue.Type:
		// passthrough
	case "":
		*fillType = solarlottery.Type
	default:
		return "", errors.Errorf(
			"%s is not a valid fill type, please use %s or %s",
			*fillType, solarlottery.Type, queue.Type)
	}

	switch types.ID(*taskType) {
	case sl.TaskTypeShift, sl.TaskTypeTicket:
		// passthrough
	case "":
		*taskType = sl.TaskTypeShift.String()
	default:
		return "", errors.Errorf(
			"%s is not a valid task type, please use %s or %s",
			*taskType, sl.TaskTypeShift, sl.TaskTypeTicket)
	}

	if *fuzz == intNoValue {
		*fuzz = 0
	}
	r, err := c.SL.MakeRotation(c.fs.Arg(0))
	if err != nil {
		return "", err
	}

	r.FillerType = types.ID(*fillType)
	r.FillSettings.Beginning = *beginning
	if r.FillSettings.Beginning.IsZero() {
		r.FillSettings.Beginning = types.NewTime(time.Now())
	}
	r.FillSettings.Period = period
	if r.FillSettings.Period.Period == "" {
		r.FillSettings.Period.Period = types.EveryWeek
	}
	r.FillSettings.Seed = *seed
	r.FillSettings.Fuzz = *fuzz

	r.TaskType = types.ID(*taskType)
	r.TaskSettings.Require.Set(sl.NeedOneAnyLevel)

	err = c.SL.AddRotation(r)
	if err != nil {
		return "", err
	}
	return "Created rotation:\n" + r.MarkdownBullets(), nil
}
