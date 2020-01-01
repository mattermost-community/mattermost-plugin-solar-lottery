// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/pkg/errors"
)

func (c *Command) leaveRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	users := ""
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.StringVar(&users, flagUsers, "", "remove other users from rotation.")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	deleted, err := c.API.LeaveRotation(users, rotation)
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", api.MarkdownUserMap(deleted))
	}

	return fmt.Sprintf("%s left rotation %s", api.MarkdownUserMap(deleted), rotation.Name), nil
}
