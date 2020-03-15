// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const (
	// commandAutopilot   = "autopilot"
	// commandForecast    = "forecast"
	// commandGuess       = "guess"
	// commandShift       = "shift"
	commandArchive     = "archive"
	commandAssign      = "assign"
	commandDebugDelete = "debug-delete"
	commandDelete      = "delete"
	commandDisqualify  = "disqualify"
	commandFinish      = "finish"
	commandGrace       = "grace"
	commandInfo        = "info"
	commandIssue       = "issue"
	commandIssueSource = "issue-source"
	commandJoin        = "join"
	commandLeave       = "leave"
	commandLimit       = "limit"
	commandList        = "list"
	commandLog         = "log"
	commandMax         = "max"
	commandMin         = "min"
	commandNew         = "new"
	commandOpen        = "open"
	commandParam       = "param"
	commandPut         = "put"
	commandQualify     = "qualify"
	commandRequire     = "require"
	commandRotation    = "rotation"
	commandShift       = "shift"
	commandShow        = "show"
	commandSkill       = "skill"
	commandStart       = "start"
	commandTask        = "task"
	commandTicket      = "ticket"
	commandFill        = "fill"
	commandSchedule    = "schedule"
	commandUnavailable = "unavailable"
	commandUser        = "user"
)

const (
	flagPFinish   = "f"
	flagPRotation = "r"
	flagPSkill    = "k"
	flagPStart    = "s"
	flagPPeriod   = "p"
)

const (
	// flagDebugRun   = "debug-run"
	// flagFillDays   = "fill-before"
	// flagNotifyDays = "notify"
	// flagOff        = "off"
	// flagSampleSize = "sample"
	flagClear    = "clear"
	flagCount    = "count"
	flagDuration = "duration"
	flagFinish   = "finish"
	flagForce    = "force"
	flagGrace    = "grace"
	flagJSON     = "json"
	flagMax      = "max"
	flagMin      = "min"
	flagPeriod   = "period"
	flagRotation = "rotation"
	flagSkill    = "skill"
	flagStart    = "start"
	flagSummary  = "summary"
)

// Command handles commands
type Command struct {
	// Config      *config.Config
	SL          sl.SL
	ConfigStore config.Store
	Context     *plugin.Context
	Args        *model.CommandArgs
	ChannelID   string

	subcommand string
	outputJson bool
	fs         *pflag.FlagSet
}

// RegisterFunc is a function that allows the runner to register commands with the mattermost server.
type RegisterFunc func(*model.Command) error

// Register should be called by the plugin to register all necessary commands
func Register(registerFunc RegisterFunc) {
	_ = registerFunc(&model.Command{
		Trigger:          constants.CommandTrigger,
		DisplayName:      "Solar Lottery",
		Description:      "team rotation scheduler",
		AutoComplete:     true,
		AutoCompleteDesc: "Schedule team rotations",
		AutoCompleteHint: fmt.Sprintf("Usage: `/%s info|rotation|shift|skill|user`.",
			constants.CommandTrigger),
	})
}

func (c *Command) commands() map[string]func([]string) (md.MD, error) {
	return map[string]func([]string) (md.MD, error){
		commandInfo:     c.info,
		commandLog:      c.log,
		commandRotation: c.rotation,
		commandSkill:    c.skill,
		commandTask:     c.task,
		commandUser:     c.user,

		"debug-clean": c.debugClean,
	}
}

// Handle should be called by the plugin when a command invocation is received from the Mattermost server.
func (c *Command) Handle() (out md.MD, err error) {
	defer func() {
		prefix := md.CodeBlock(c.Args.Command) + "\n"
		if err != nil {
			prefix += md.Markdownf("Command failed. Error: **%s**\n", err.Error())
		}
		out = prefix + out
	}()

	command, parameters, err := c.validate()
	if err != nil {
		return "", err
	}
	c.subcommand = command
	return c.handleCommand(c.commands(), parameters)
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
	if command != "/"+constants.CommandTrigger {
		return "", nil, errors.Errorf("%q is not a supported command and should not have been invoked. Please contact your system administrator", command)
	}

	return command, split[1:], nil
}

func (c *Command) handleCommand(
	subcommands map[string]func([]string) (md.MD, error),
	parameters []string,
) (md.MD, error) {
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

func (c *Command) normalOut(out md.Markdowner, err error) (md.MD, error) {
	if err != nil {
		return "", err
	}
	if c.outputJson {
		out = md.JSONBlock(out)
	}
	return out.Markdown(), nil
}

func (c *Command) flagUsage() md.MD {
	usage := c.subcommand
	if c.fs != nil {
		usage += " [flags...]\n\nFlags:\n" + c.fs.FlagUsages()
	}
	return md.Markdownf("Usage:\n%s", md.CodeBlock(usage))
}

func (c *Command) subUsage(subcommands map[string]func([]string) (md.MD, error)) md.MD {
	subs := []string{}
	for sub := range subcommands {
		subs = append(subs, sub)
	}
	sort.Strings(subs)
	usage := fmt.Sprintf("`%s %s`", c.subcommand, strings.Join(subs, "|"))
	return md.Markdownf("Usage: %s\nUse `%s <subcommand> help` for more info.",
		usage, c.subcommand)
}

func (c *Command) debugClean(parameters []string) (md.MD, error) {
	return "Cleaned the KV store", c.SL.Clean()
}

func (c *Command) assureFS() *pflag.FlagSet {
	if c.fs == nil {
		c.fs = pflag.NewFlagSet("", pflag.ContinueOnError)
		c.fs.BoolVar(&c.outputJson, flagJSON, false, "output as JSON")
	}
	return c.fs
}

func (c *Command) withFlagRotation() {
	c.assureFS().StringP(flagRotation, flagPRotation, "", "rotation reference")
}

func (c *Command) withFlagStartFinish(actingUser *sl.User) (*types.Time, *types.Time) {
	start := actingUser.Time(types.NewTime())
	finish := start
	c.assureFS().VarP(&start, flagStart, flagPStart, "start of the interval")
	c.assureFS().VarP(&finish, flagFinish, flagPFinish, "end of the interval")
	return &start, &finish
}

func (c *Command) withFlagStart(actingUser *sl.User) *types.Time {
	start := actingUser.Time(types.NewTime())
	c.assureFS().VarP(&start, flagStart, flagPStart, "start time")
	return &start
}

func (c *Command) withFlagPeriod() *types.Period {
	p := types.Period{}
	c.assureFS().VarP(&p, flagPeriod, flagPPeriod, "recurrence period")
	return &p
}

func (c *Command) withFlagClear() *bool {
	return c.assureFS().Bool(flagClear, false, "mark as available by clearing all overlapping unavailability events")
}

func (c *Command) withFlagMin() *bool {
	return c.assureFS().Bool(flagMin, false, "add/update a minimum headcount requirement for a skill level")
}

func (c *Command) withFlagMax() *bool {
	return c.assureFS().Bool(flagMax, false, "add/update a maximum headcount requirement for a skill level")
}

func (c *Command) withFlagDuration() *time.Duration {
	return c.assureFS().Duration(flagDuration, 0, "add a grace period to users after finishing a task")
}

func (c *Command) withFlagCount() *int {
	return c.assureFS().Int(flagCount, 1, "number of users")
}

func (c *Command) withFlagSkill() *string {
	return c.assureFS().StringP(flagSkill, flagPSkill, "", "Skill name")
}

func (c *Command) withFlagSkillLevel() *sl.SkillLevel {
	var skillLevel sl.SkillLevel
	c.assureFS().VarP(&skillLevel, flagSkill, flagPSkill, "Skill-level")
	return &skillLevel
}

func (c *Command) withFlagSummary() *string {
	return c.assureFS().String(flagSummary, "", "task summary")
}

func (c *Command) withFlagForce() *bool {
	return c.assureFS().Bool(flagForce, false, "ignore constraints")
}

func (c *Command) resolveUsernames(args []string) (mattermostUserIDs *types.IDSet, err error) {
	mattermostUserIDs = types.NewIDSet()
	// if no args provided, return the acting user
	if len(args) == 0 {
		user, err := c.SL.ActingUser()
		if err != nil {
			return nil, err
		}
		mattermostUserIDs.Set(user.MattermostUserID)
		return mattermostUserIDs, nil
	}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "@") {
			return nil, errors.New("`@username`'s expected")
		}
		arg = arg[1:]
		user, err := c.SL.LoadMattermostUsername(arg)
		if err != nil {
			return nil, err
		}
		mattermostUserIDs.Set(user.MattermostUserID)
	}

	return mattermostUserIDs, nil
}

func (c *Command) resolveRotationUsernames() (types.ID, *types.IDSet, error) {
	ref, _ := c.fs.GetString(flagRotation)
	usernames := []string{}
	rotationID := types.ID(ref)

	for _, arg := range c.fs.Args() {
		if strings.HasPrefix(arg, "@") {
			usernames = append(usernames, arg)
		} else {
			if rotationID != "" {
				return "", nil, errors.Errorf("rotation %s is already specified, cant't interpret %s", rotationID, arg)
			}
			rotationID = types.ID(arg)
		}
	}

	var err error
	if rotationID == "" {
		return "", nil, errors.New("rotation must be specified")
	}
	// explicit ref is used as is
	if ref == "" {
		rotationID, err = c.SL.ResolveRotationName(string(rotationID))
		if err != nil {
			return "", nil, err
		}
	}

	mattermostUserIDs, err := c.resolveUsernames(usernames)
	if err != nil {
		return "", nil, err
	}
	return rotationID, mattermostUserIDs, nil
}

func (c *Command) resolveTaskIDUsernames() (types.ID, *types.IDSet, error) {
	args := c.fs.Args()
	if len(args) == 0 {
		return "", nil, errors.New("Task ID is required")
	}
	usernames := []string{}
	taskID := types.ID(args[0])
	args = args[1:]
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			usernames = append(usernames, arg)
		} else {
			return "", nil, errors.Errorf("Unexpected argument: %s, expected @usernames", arg)
		}
	}

	mattermostUserIDs, err := c.resolveUsernames(usernames)
	if err != nil {
		return "", nil, err
	}
	return taskID, mattermostUserIDs, nil
}

func (c *Command) resolveRotation() (types.ID, error) {
	var err error
	ref, _ := c.fs.GetString(flagRotation)
	rotationID := types.ID(ref)
	if ref == "" {
		if len(c.fs.Args()) < 1 {
			return "", errors.New("no rotation specified")
		}
		rotationID, err = c.SL.ResolveRotationName(c.fs.Arg(0))
		if err != nil {
			return "", err
		}
	}
	return rotationID, nil
}
