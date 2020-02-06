// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c *Command) listRotations(parameters []string) (string, error) {
	if len(parameters) > 0 {
		return c.subUsage(nil), errors.New("unexpected parameters")
	}
	rotations, err := c.SL.LoadKnownRotations()
	if err != nil {
		return "", err
	}
	if len(rotations) == 0 {
		return "*none*", nil
	}

	out := ""
	for id := range rotations {
		out += fmt.Sprintf("- %s\n", id)
	}
	return out, nil
}
