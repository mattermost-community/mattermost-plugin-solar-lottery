// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) showRotation(parameters []string) (string, error) {
	fs := newRotationUsersFlagSet()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs.FlagSet), err
	}

	r, _, err := c.rotationUsers(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.ExpandRotation(r)
	if err != nil {
		return "", err
	}
	return r.MarkdownBullets(), nil
}
