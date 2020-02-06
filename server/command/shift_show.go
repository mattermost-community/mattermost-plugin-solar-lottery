// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

func (c *Command) showShift(parameters []string) (string, error) {
	return c.doShift(
		parameters,
		nil,
		func(fs *pflag.FlagSet, rotation *sl.Rotation, shiftNumber int) (string, error) {
			shift, err := c.SL.OpenShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shift), nil
		})
}
