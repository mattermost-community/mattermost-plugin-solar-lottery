// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationNew(parameters []string) (md.MD, error) {
	err := c.parse(parameters)
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
	// r.Duration = 24 * time.Hour
	r.TaskSettings.Require.Set(sl.NeedOneAnyLevel)

	err = c.SL.AddRotation(r)
	if err != nil {
		return "", err
	}

	return "Created rotation:\n" + r.MarkdownBullets(), nil
}
