// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationNew(parameters []string) (md.MD, error) {
	fs := c.assureFS()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	if fs.Arg(0) == "" {
		return c.flagUsage(), errors.Errorf("must specify rotation name")
	}

	r, err := c.SL.MakeRotation(fs.Arg(0))
	if err != nil {
		return "", err
	}

	// TODO parameterize rotation defaults
	r.TaskFillerType = solarlottery.Type
	r.TaskMaker.Type = sl.TicketMaker
	r.TaskMaker.TicketDefaultDuration = 24 * time.Hour
	r.TaskMaker.Require.Set(sl.NeedOneAnyLevel)

	err = c.SL.AddRotation(r)
	if err != nil {
		return "", err
	}

	return "Created rotation:\n" + r.MarkdownBullets(), nil
}
