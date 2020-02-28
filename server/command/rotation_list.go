// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/pkg/errors"
)

func (c *Command) listRotations(parameters []string) (string, error) {
	fs := newFS()
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if len(fs.Args()) > 0 {
		return c.subUsage(nil), errors.New("unexpected parameters")
	}

	rotations, err := c.SL.LoadActiveRotations()
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(rotations), nil
	}
	if rotations.Len() == 0 {
		return "*none*", nil
	}
	out := ""
	for _, id := range rotations.IDs() {
		out += fmt.Sprintf("- %s\n", id)
	}
	return out, nil
}
