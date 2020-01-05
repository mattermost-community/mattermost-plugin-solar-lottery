// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func withRotationAddFlags(fs *pflag.FlagSet, start *string, period *api.Period) {
	fs.StringVarP(start, flagStart, flagPStart, "", fmt.Sprintf("rotation start date formatted as %s. It must be provided at creation and **can not be modified** later.", api.DateFormat))
	fs.Var(period, flagPeriod, "rotation period 1w, 2w, or 1m")
}

func (c *Command) addRotation(parameters []string) (string, error) {
	var rotationName, start string
	var period api.Period
	var size, grace int
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	withRotationAddFlags(fs, &start, &period)
	withRotationUpdateFlags(fs, &size, &grace)
	fs.StringVarP(&rotationName, flagRotation, flagPRotation, "", "specify rotation name")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	if rotationName == "" {
		return c.flagUsage(fs), errors.Errorf("must specify rotation name, use `--%s`", flagRotation)
	}

	rotation, err := c.API.MakeRotation(rotationName)
	if err != nil {
		return "", err
	}
	rotation.Period = period.String()
	rotation.Start = start
	rotation.Size = size
	rotation.Grace = grace

	err = c.API.AddRotation(rotation)
	if err != nil {
		return "", err
	}

	return "Created rotation:\n" + c.API.MarkdownRotationBullets(rotation), nil
}
