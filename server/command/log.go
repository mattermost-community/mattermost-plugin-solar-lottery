// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) log(parameters []string) (md.MD, error) {
	var level string
	var verbose bool
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&level, "level", "info", "log level")
	fs.BoolVar(&verbose, "context", false, "include context with log messages")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	sc := c.SL.Config().StoredConfig
	sc.AdminLogLevel = level
	sc.AdminLogVerbose = verbose
	c.ConfigStore.SaveConfig(sc)
	return "Dispatched config update.", nil
}
