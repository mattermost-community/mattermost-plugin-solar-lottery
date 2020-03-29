// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
)

func (c *Command) leaveShift(parameters []string) (string, error) {
	var usernames string
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to leave the shift.")
		},
		func(fs *pflag.FlagSet, rotation *sl.Rotation, shiftNumber int) (string, error) {
			_, deleted, err := c.SL.LeaveShift(usernames, rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s left %s", deleted.MarkdownWithSkills(), rotation.ShiftRef(shiftNumber)), nil
		})
}
