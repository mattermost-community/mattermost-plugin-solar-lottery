// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) listShifts(parameters []string) (string, error) {
	numShifts := 0
	return c.doShift(parameters,
		func(fs *pflag.FlagSet) {
			fs.IntVarP(&numShifts, flagNumber, flagPNumber, 3, "Number of shifts to list")
		},
		func(fs *pflag.FlagSet, rotation *sl.Rotation, shiftNumber int) (string, error) {
			shifts, err := c.SL.ListShifts(rotation, shiftNumber, numShifts)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shifts), nil
		})
}
