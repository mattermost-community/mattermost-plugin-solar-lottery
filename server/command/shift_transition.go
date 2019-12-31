// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) commitShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, start int) (string, error) {
			err := c.API.CommitShift(rotation, start)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Committed shift #%v", start), nil
		})
}

func (c *Command) finishShift(parameters []string) (string, error) {
	return c.doShift(parameters,
		nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.FinishShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished shift #%v", shiftNumber), nil
		})
}

func (c *Command) startShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.StartShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Finished shift #%v", shiftNumber), nil
		})
}

func (c *Command) debugDeleteShift(parameters []string) (string, error) {
	return c.doShift(parameters, nil,
		func(rotation *api.Rotation, shiftNumber int) (string, error) {
			err := c.API.DebugDeleteShift(rotation, shiftNumber)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Deleted shift #%v", shiftNumber), nil
		})
}
