// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) fillShift(parameters []string) (string, error) {
	return c.doShift(parameters,
		nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, ready, _, err := c.API.IsShiftReady(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			if ready {
				return fmt.Sprintf("%s is ready, no fill required.\nUsers: %s.",
					rotation.ShiftRef(shiftNumber),
					rotation.ShiftUsers(shift).MarkdownWithSkills()), nil
			}

			shift, users, err := c.API.FillShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s is filled, added: %s.",
				rotation.ShiftRef(shiftNumber), users.MarkdownWithSkills()), nil
		})
}
