// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/pkg/errors"
)

func (c *Command) rotationList(parameters []string) (md.MD, error) {
	err := c.parse(parameters)
	if len(c.fs.Args()) > 0 {
		return c.subUsage(nil), errors.New("unexpected parameters")
	}

	rotations, err := c.SL.LoadActiveRotations()
	if err != nil {
		return "", err
	}

	if c.outputJson {
		return md.JSONBlock(rotations), nil
	}
	if rotations.Len() == 0 {
		return "*none*", nil
	}
	out := md.MD("")
	for _, id := range rotations.IDs() {
		out += md.Markdownf("- %s\n", id)
	}
	return out, nil
}
