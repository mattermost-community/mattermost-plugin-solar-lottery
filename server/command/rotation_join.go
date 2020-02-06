// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
)

func (c *Command) joinRotation(parameters []string) (string, error) {
	var rotationID, rotationName, start string
	users := ""
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.StringVarP(&users, flagUsers, flagPUsers, "", "add other users to rotation.")
	fs.StringVarP(&start, flagStart, flagPStart, "", fmt.Sprintf("date for user to start, e.g. %s.", sl.DateFormat))
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.SL.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	starting := time.Now()
	if start != "" {
		starting, err = time.Parse(sl.DateFormat, start)
		if err != nil {
			return c.flagUsage(fs), err
		}
	}
	added, err := c.SL.JoinRotation(users, rotation, starting)
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", added.Markdown())
	}

	return fmt.Sprintf("%s joined rotation %s", added.Markdown(), rotation.Name), nil
}
