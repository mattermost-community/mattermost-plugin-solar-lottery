// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/spf13/pflag"
)

func (c *Command) fillShift(parameters []string) (string, error) {
	var autofill bool
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&autofill, flagAutofill, false, "autofill shift if needed")
		},
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, ready, whyNot, before, added, err := c.API.FillShift(rotation, shiftNumber, autofill)
			switch {
			case err != nil:
				return "", err

			case !ready:
				return fmt.Sprintf("%s is not ready.\nUsers: %s.\nReason: %s.",
					api.MarkdownShift(rotation, shiftNumber, shift),
					api.MarkdownUserMapWithSkills(shift.Users),
					whyNot), nil

			case ready && len(added) == 0:
				return fmt.Sprintf("%s is ready, no fill required.\nUsers: %s.",
					api.MarkdownShift(rotation, shiftNumber, shift),
					api.MarkdownUserMapWithSkills(shift.Users)), nil

			case ready && len(before) > 0:
				return fmt.Sprintf("%s is ready.\nUsers already in the shift: %s\nAdded users: %s.",
					api.MarkdownShift(rotation, shiftNumber, shift),
					api.MarkdownUserMapWithSkills(before),
					api.MarkdownUserMapWithSkills(added)), nil

			default:
				return fmt.Sprintf("%s is ready.\nAdded users: %s.",
					api.MarkdownShift(rotation, shiftNumber, shift),
					api.MarkdownUserMapWithSkills(added)), nil

			}
		})
}
