// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) showRotation(parameters []string) (string, error) {
	fs := newRotationFS()
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	r, _, err := c.rotationUsers(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.ExpandRotation(r)
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return r.MarkdownBullets(), nil
}
