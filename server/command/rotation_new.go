// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) rotationNew(parameters []string) (md.MD, error) {
	seed := c.withFlagSeed()
	begin, err := c.withFlagBeginning()
	if err != nil {
		return "", err
	}
	err = c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	if c.fs.Arg(0) == "" {
		return c.flagUsage(), errors.Errorf("must specify rotation name")
	}

	r, err := c.SL.MakeRotation(c.fs.Arg(0))
	if err != nil {
		return "", err
	}

	// TODO parameterize rotation defaults
	r.FillerType = solarlottery.Type
	r.TaskType = sl.TaskTypeTicket
	r.TaskSettings.Require.Set(sl.NeedOneAnyLevel)
	r.Seed = *seed

	if begin.IsZero() {
		r.Beginning = types.NewTime(time.Now())
	} else {
		r.Beginning = *begin
	}

	err = c.SL.AddRotation(r)
	if err != nil {
		return "", err
	}

	return "Created rotation:\n" + r.MarkdownBullets(), nil
}
