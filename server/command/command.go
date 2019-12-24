// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils"
)

const (
	commandAdd           = "add"
	commandArchive       = "archive"
	commandCommit        = "commit"
	commandDebugDelete   = "debug-delete"
	commandDebugRotation = "debug-rotation"
	commandDebugShift    = "debug-shift"
	commandDelete        = "delete"
	commandDisqualify    = "disqualify"
	commandFill          = "fill"
	commandFinish        = "finish"
	commandForecast      = "forecast"
	commandHeat          = "heat"
	commandHelp          = "help"
	commandInfo          = "info"
	commandJoin          = "join"
	commandLeave         = "leave"
	commandList          = "list"
	commandNeed          = "need"
	commandQualify       = "qualify"
	commandRotation      = "rotation"
	commandSchedule      = "schedule"
	commandShift         = "shift"
	commandShow          = "show"
	commandSkill         = "skill"
	commandStart         = "start"
	commandUnavailable   = "unavailable"
	commandUpdate        = "update"
	commandUser          = "user"
)

const (
	flagPEnd      = "e"
	flagPLevel    = "l"
	flagPNumber   = "n"
	flagPRotation = "r"
	flagPSkill    = "k"
	flagPStart    = "s"
	flagPUsers    = "u"
)

const (
	flagAutofill   = "autofill"
	flagShifts     = "shifts"
	flagDeleteNeed = "delete-need"
	flagEnd        = "end"
	flagGrace      = "grace"
	flagLevel      = "level"
	flagMax        = "max"
	flagMin        = "min"
	flagNumber     = "number"
	flagPadding    = "padding"
	flagPeriod     = "period"
	flagRotation   = "rotation"
	flagRotationID = "rotation-id"
	flagSize       = "size"
	flagSkill      = "skill"
	flagStart      = "start"
	flagUsers      = "users"
)

// Command handles commands
type Command struct {
	Context   *plugin.Context
	Args      *model.CommandArgs
	ChannelID string
	Config    *config.Config
	API       api.API
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
		AutoCompleteHint: fmt.Sprintf("some commands: `%s,%s,%s,%s,%s`. Type `/%s help` for more information.",
			commandShift, commandForecast, commandJoin, commandLeave, commandUser, config.CommandTrigger),
	})
}

func (c *Command) handleCommand(subcommands map[string]func([]string) (string, error),
	parameters []string, usage string) (string, error) {
	if len(parameters) == 0 {
		return usage, errors.New("expected a (sub-)command")
	}

	f := subcommands[parameters[0]]
	if f == nil {
		return usage, errors.Errorf("unknown command: %s", parameters[0])
	}

	return f(parameters[1:])
}

// Handle should be called by the plugin when a command invocation is received from the Mattermost server.
func (c *Command) Handle() (out string, err error) {
	subcommands := map[string]func([]string) (string, error){
		commandAdd:      c.add,
		commandNeed:     c.need,
		commandDelete:   c.delete,
		commandHelp:     c.help,
		commandInfo:     c.info,
		commandRotation: c.rotation,
		commandShift:    c.shift,
		commandShow:     c.show,
		commandUser:     c.user,
		commandList:     c.list,
		commandJoin:     c.join,
		commandLeave:    c.leave,
		commandForecast: c.forecast,
	}

	defer func() {
		prefix := utils.CodeBlock(c.Args.Command) + "\n"
		if err != nil {
			prefix += "Command failed. Error: **" + err.Error() + "**\n"
		}
		out = prefix + out
	}()

	parameters, err := c.validate()
	if err != nil {
		return "", err
	}

	return c.handleCommand(subcommands, parameters, "TODO Usage: /slottery do stuff")
}

func (c *Command) validate() ([]string, error) {
	if c.Context == nil || c.Args == nil {
		return nil, errors.New("invalid arguments to command.Handler. Please contact your system administrator")
	}
	split := strings.Fields(c.Args.Command)
	if len(split) < 2 {
		return nil, errors.New("io subcommand specify, nothing to do")
	}
	command := split[0]
	if command != "/"+config.CommandTrigger {
		return nil, errors.Errorf("%q is not a supported command and should not have been invoked. Please contact your system administrator", command)
	}

	return split[1:], nil
}
