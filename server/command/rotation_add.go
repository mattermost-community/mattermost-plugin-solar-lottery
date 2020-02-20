// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func (c *Command) addRotation(parameters []string) (string, error) {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	if fs.Arg(0) == "" {
		return c.flagUsage(fs), errors.Errorf("must specify rotation name")
	}

	r, err := c.SL.MakeRotation(fs.Arg(0))
	if err != nil {
		return "", err
	}
	err = c.SL.AddRotation(r)
	if err != nil {
		return "", err
	}

	return "Created rotation:\n" + r.MarkdownBullets(), nil
}
