// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/pkg/errors"
)

func (c *Command) joinRotation(parameters []string) (string, error) {
	var rotationID, rotationName, start string
	users := ""
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.StringVarP(&users, flagUsers, flagPUsers, "", "add nother users to rotation.")
	fs.StringVarP(&start, flagStart, flagPStart, "", fmt.Sprintf("date for user to start, e.g. %s.", api.DateFormat))
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

	starting := time.Now()
	if start != "" {
		starting, err = time.Parse(api.DateFormat, start)
		if err != nil {
			return c.flagUsage(fs), err
		}
	}
	added, err := c.API.JoinRotation(users, rotation, starting)
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", c.API.MarkdownUsers(added))
	}

	return fmt.Sprintf("%s joined rotation %s", c.API.MarkdownUsers(added), rotation.Name), nil
}
