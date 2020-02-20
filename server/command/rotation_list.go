// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
)

func (c *Command) listRotations(parameters []string) (string, error) {
	if len(parameters) > 0 {
		return c.subUsage(nil), errors.New("unexpected parameters")
	}
	rotations, err := c.SL.LoadActiveRotations()
	if err != nil {
		return "", err
	}
	if rotations.Len() == 0 {
		return "*none*", nil
	}

	out := ""
	rotations.ForEach(func(id string) {
		out += fmt.Sprintf("- %s\n", kvstore.NameFromID(id))
	})
	return out, nil
}
