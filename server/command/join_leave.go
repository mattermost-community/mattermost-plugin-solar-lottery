// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c *Command) join(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]
	err := c.API.JoinRotation(rotationName, 0)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Joined rotation %s", rotationName), nil
}

func (c *Command) leave(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]
	err := c.API.LeaveRotation(rotationName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Left rotation %s", rotationName), nil
}
