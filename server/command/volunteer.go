// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/spf13/pflag"
)

func (c *Command) volunteer(parameters ...string) (string, error) {
	var usernames string
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.StringVar(&usernames, "users", "", "add nother users to rotation.")
		},
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.VolunteerUsers(usernames, rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Volunteered %s for shift #%v in %s", usernames, shiftNumber, api.MarkdownRotation(rotation)), nil
		})
}
