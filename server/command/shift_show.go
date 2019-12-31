// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
	"github.com/spf13/pflag"
)

func (c *Command) showShift(parameters []string) (string, error) {
	return c.doShift(
		parameters,
		nil,
		func(fs *pflag.FlagSet, rotation *api.Rotation, shiftNumber int) (string, error) {
			shift, err := c.API.OpenShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return utils.JSONBlock(shift), nil
		})
}
