// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) newRotation(parameters []string) (md.MD, error) {
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
	err = c.SL.AddRotation(r)
	if err != nil {
		return "", err
	}

	return "Created rotation:\n" + r.MarkdownBullets(), nil
}
