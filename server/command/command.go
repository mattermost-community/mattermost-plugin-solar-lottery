// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

const (
	commandAdd         = "add"
	commandArchive     = "archive"
	commandAutopilot   = "autopilot"
	commandDebugDelete = "debug-delete"
	commandDelete      = "delete"
	commandDisqualify  = "disqualify"
	commandFill        = "fill"
	commandFinish      = "finish"
	commandForecast    = "forecast"
	commandGuess       = "guess"
	commandInfo        = "info"
	commandJoin        = "join"
	commandLeave       = "leave"
	commandList        = "list"
	commandLog         = "log"
	commandNeed        = "need"
	commandOpen        = "open"
	commandQualify     = "qualify"
	commandRotation    = "rotation"
	commandShift       = "shift"
	commandShow        = "show"
	commandSkill       = "skill"
	commandStart       = "start"
	commandUnavailable = "unavailable"
	commandUpdate      = "update"
	commandUser        = "user"
)

const (
	flagPEnd      = "e"
	flagPLevel    = "l"
	flagPNumber   = "n"
	flagPRotation = "r"
	flagPShift    = "s"
	flagPSkill    = "k"
	flagPStart    = "s"
	flagPUsers    = "u"
)

const (
	flagClear      = "clear"
	flagDebugRun   = "debug-run"
	flagDeleteNeed = "delete-need"
	flagEnd        = "end"
	flagFill       = "fill"
	flagFillDays   = "fill-before"
	flagGrace      = "grace"
	flagJSON       = "json"
	flagLevel      = "level"
	flagMax        = "max"
	flagMin        = "min"
	flagNotifyDays = "notify"
	flagNumber     = "number"
	flagOff        = "off"
	flagPeriod     = "period"
	flagRotation   = "rotation"
	flagRotationID = "rotation-id"
	flagSampleSize = "sample"
	flagShift      = "shift"
	flagSize       = "size"
	flagSkill      = "skill"
	flagStart      = "start"
	flagType       = "type"
	flagUsers      = "users"
)

// Command handles commands
type Command struct {
	Context   *plugin.Context
	Args      *model.CommandArgs
	ChannelID string
	Config    *config.Config
	API       api.API

	subcommand string
}

// RegisterFunc is a function that allows the runner to register commands with the mattermost server.
type RegisterFunc func(*model.Command) error

// Register should be called by the plugin to register all necessary commands
func Register(registerFunc RegisterFunc) {
	_ = registerFunc(&model.Command{
		Trigger:          config.CommandTrigger,
		DisplayName:      "Solar Lottery",
		Description:      "team rotation scheduler",
		AutoComplete:     true,
		AutoCompleteDesc: "Schedule team rotations",
		AutoCompleteHint: fmt.Sprintf("Usage: `/%s info|rotation|shift|skill|user`.",
			config.CommandTrigger),
	})
}

// Handle should be called by the plugin when a command invocation is received from the Mattermost server.
func (c *Command) Handle() (out string, err error) {
	subcommands := map[string]func([]string) (string, error){
		commandInfo:     c.info,
		commandRotation: c.rotation,
		commandShift:    c.shift,
		commandSkill:    c.skill,
		commandUser:     c.user,
		commandLog:      c.log,

		"debug-clean": c.debugClean,
	}

	defer func() {
		prefix := utils.CodeBlock(c.Args.Command) + "\n"
		if err != nil {
			prefix += "Command failed. Error: **" + err.Error() + "**\n"
		}
		out = prefix + out
	}()

	command, parameters, err := c.validate()
	if err != nil {
		return "", err
	}
	c.subcommand = command
	return c.handleCommand(subcommands, parameters)
}

func (c *Command) validate() (string, []string, error) {
	if c.Context == nil || c.Args == nil {
		return "", nil, errors.New("invalid arguments to command.Handler. Please contact your system administrator")
	}
	split := strings.Fields(c.Args.Command)
	if len(split) < 2 {
		return "", nil, errors.New("no subcommand specified, nothing to do")
	}
	command := split[0]
	if command != "/"+config.CommandTrigger {
		return "", nil, errors.Errorf("%q is not a supported command and should not have been invoked. Please contact your system administrator", command)
	}

	return command, split[1:], nil
}

func (c *Command) handleCommand(subcommands map[string]func([]string) (string, error),
	parameters []string) (string, error) {
	if len(parameters) == 0 {
		return c.subUsage(subcommands), errors.New("expected a (sub-)command")
	}

	f := subcommands[parameters[0]]
	if f == nil {
		return c.subUsage(subcommands), errors.Errorf("unknown command: %s", parameters[0])
	}
	c.subcommand += " " + parameters[0]

	return f(parameters[1:])
}

func (c *Command) flagUsage(fs *pflag.FlagSet) string {
	usage := c.subcommand
	if fs != nil {
		usage += " [flags...]\n\nFlags:\n" + fs.FlagUsages()
	}
	return fmt.Sprintf("Usage:\n" + utils.CodeBlock(usage))
}

func (c *Command) subUsage(subcommands map[string]func([]string) (string, error)) string {
	subs := []string{}
	for sub := range subcommands {
		subs = append(subs, sub)
	}
	sort.Strings(subs)
	usage := fmt.Sprintf("`%s %s`", c.subcommand, strings.Join(subs, "|"))
	return fmt.Sprintf("Usage: %s\nUse `%s <subcommand> help` for more info.",
		usage, c.subcommand)
}

func (c *Command) debugClean(parameters []string) (string, error) {
	return "Cleaned the KV store", c.API.Clean()
}
