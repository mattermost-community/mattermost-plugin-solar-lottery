// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

func (c *Command) showRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := newRotationFlagSet(&rotationID, &rotationName)
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
	err = c.SL.ExpandRotation(rotation)
	if err != nil {
		return "", err
	}
	return rotation.MarkdownBullets(), nil
}
