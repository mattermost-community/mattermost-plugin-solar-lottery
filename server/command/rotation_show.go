// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) showRotation(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err := c.resolveRotation(fs)
	if err != nil {
		return "", err
	}
	r, err := c.SL.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return r.MarkdownBullets(), nil
}
