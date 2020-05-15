// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/spf13/pflag"
)

func (c *Command) main(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"info":     c.info,
		"rotation": c.rotation,
		"skill":    c.skill,
		"task":     c.task,
		"user":     c.user,

		"debug-log":   c.debugLog,
		"debug-clean": c.debugClean,
	}

	return c.run(subcommands, parameters)
}

func (c *Command) rotation(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"archive":      c.rotationArchive,
		"autopilot":    c.rotationAutopilot,
		"debug-delete": c.rotationDebugDelete,
		"list":         c.rotationList,
		"new":          c.rotationNew,
		"set":          c.rotationSet,
		"show":         c.rotationShow,
	}
	return c.run(subcommands, parameters)
}

func (c *Command) rotationSet(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"autopilot": c.rotationSetAutopilot,
		"fill":      c.rotationSetFill,
		"limit":     c.rotationSetLimit,
		"require":   c.rotationSetRequire,
		"task":      c.rotationSetTask,
	}
	return c.run(subcommands, parameters)
}

func (c *Command) task(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"assign":   c.taskAssign,
		"unassign": c.taskUnassign,
		"fill":     c.taskFill,
		"schedule": c.taskTransition(sl.TaskStateScheduled),
		"start":    c.taskTransition(sl.TaskStateStarted),
		"finish":   c.taskTransition(sl.TaskStateFinished),
		"new":      c.taskNew,
		"show":     c.taskShow,
	}
	return c.run(subcommands, parameters)
}

func (c *Command) taskNew(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"ticket": c.taskNewTicket,
		"shift":  c.taskNewShift,
	}

	return c.run(subcommands, parameters)
}

func (c *Command) user(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"disqualify":  c.userDisqualify,
		"qualify":     c.userQualify,
		"show":        c.userShow,
		"unavailable": c.userUnavailable,
		"join":        c.userJoin,
		"leave":       c.userLeave,
	}
	return c.run(subcommands, parameters)
}

func (c *Command) skill(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		"new":    c.skillNew,
		"delete": c.skillDelete,
		"list":   c.skillList,
	}
	return c.run(subcommands, parameters)
}

func (c *Command) debugClean(parameters []string) (md.MD, error) {
	return "Cleaned the KV store", c.SL.Clean()
}

func (c *Command) debugLog(parameters []string) (md.MD, error) {
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
