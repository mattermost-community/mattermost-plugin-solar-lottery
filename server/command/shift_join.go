// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
)

func (c *Command) joinShift(parameters []string) (string, error) {
	var usernames string
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to join the shift.")
		},
		func(fs *pflag.FlagSet, rotation *sl.Rotation, shiftNumber int) (string, error) {
			_, joined, err := c.SL.JoinShift(usernames, rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s joined %s", joined.MarkdownWithSkills(), rotation.ShiftRef(shiftNumber)), nil
		})
}
