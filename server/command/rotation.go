// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (c *Command) rotation(parameters ...string) (string, error) {
	subcommands := map[string]func(...string) (string, error){
		"archive":  c.archiveRotation,
		"create":   c.createRotation,
		"forecast": c.forecast,
		"list":     c.listRotations,
		"need":     c.updateRotationNeed,
		"show":     c.showRotation,
		"update":   c.updateRotation,
	}
	errUsage := errors.Errorf("Invalid subcommand. Usage:\n"+
		"- `%s rotation list` - list all known rotations\n"+
		"- `%s rotation archive|create|forecast|show|update] <rotation-name>` - view or modify specific rotations\n"+
		"\n"+
		"Use `%s rotation subcommand --help` for more information.\n",
		config.CommandTrigger, config.CommandTrigger, config.CommandTrigger)

	if len(parameters) == 0 {
		return "", errUsage
	}

	f := subcommands[parameters[0]]
	if f == nil {
		return "", errUsage
	}

	return f(parameters[1:]...)
}

func (c *Command) showRotation(parameters ...string) (string, error) {
	fs := flag.NewFlagSet("showRotation", flag.ContinueOnError)

	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "rotation show <rotation-name>")
	if err != nil {
		return "", err
	}
	return api.MarkdownRotationWithDetails(rotation), nil
}

func (c *Command) listRotations(parameters ...string) (string, error) {
	if len(parameters) > 0 {
		return "", errors.Errorf(commandUsage("rotation forecast <rotation-name>", nil))
	}
	rotations, err := c.API.LoadKnownRotations()
	if err != nil {
		return "", err
	}
	if len(rotations) == 0 {
		return "*none*", nil
	}

	out := ""
	for id := range rotations {
		out += fmt.Sprintf("- %s\n", id)
	}
	return out, nil
}

func (c *Command) archiveRotation(parameters ...string) (string, error) {
	fs := flag.NewFlagSet("archiveRotation", flag.ContinueOnError)

	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "rotation archive <rotation-name>")
	if err != nil {
		return "", err
	}

	err = c.API.ArchiveRotation(rotation)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", rotation.Name)
	}

	return "Deleted rotation " + rotation.Name, nil
}

func (c *Command) createRotation(parameters ...string) (string, error) {
	var start string
	var period api.Period
	var size, paddingWeeks int
	fs := flag.NewFlagSet("createRotation", flag.ContinueOnError)
	withRotationCreate(fs, &start, &period)
	withRotationOptions(fs, &size, &paddingWeeks)
	usage := commandUsage("rotation create <rotation-name>", fs)
	err := fs.Parse(parameters)
	if err != nil {
		return "", errors.Errorf("**%s**.\n\n%s", err.Error(), usage)
	}
	if len(fs.Args()) != 1 || start == "" || period.String() == "" {
		return "", errors.New(usage)
	}

	rotation, err := c.API.MakeRotation(fs.Arg(0))
	if err != nil {
		return "", err
	}
	rotation.Period = period.String()
	rotation.Start = start
	rotation.Size = size
	rotation.PaddingWeeks = paddingWeeks

	err = c.API.AddRotation(rotation)
	if err != nil {
		return "", err
	}

	return "Created rotation " + api.MarkdownRotationWithDetails(rotation), nil
}

func (c *Command) updateRotation(parameters ...string) (string, error) {
	var size, paddingWeeks int
	fs := flag.NewFlagSet("updateRotation", flag.ContinueOnError)
	withRotationOptions(fs, &size, &paddingWeeks)

	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "rotation update <rotation-name>")
	if err != nil {
		return "", err
	}

	err = c.API.UpdateRotation(rotation, func(rotation *api.Rotation) error {
		if paddingWeeks != 0 {
			rotation.PaddingWeeks = paddingWeeks
		}
		if size != 0 {
			rotation.Size = size
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return "Updated rotation " + api.MarkdownRotationWithDetails(rotation), nil
}

func (c *Command) updateRotationNeed(parameters ...string) (string, error) {
	var name, skill, level, removeName string
	var min, max int
	fs := flag.NewFlagSet("updateRotationNeed", flag.ContinueOnError)
	withRotationNeed(fs, &name, &skill, &level, &min, &max, &removeName)

	rotation, err := c.parseRotationFlagsAndLoad(fs, parameters, "rotation need <rotation-name>")
	if err != nil {
		return "", err
	}
	if len(name) == 0 && len(removeName) == 0 {
		return "", errors.Errorf("only one of --need and --remove-need must be specified. %s",
			commandUsage("rotation need <rotation-name>", fs))
	}
	if len(name) != 0 && len(removeName) != 0 {
		return "", errors.Errorf("one of --need and --remove-need must be specified. %s",
			commandUsage("rotation need <rotation-name>", fs))
	}

	var updatef func(rotation *api.Rotation) error
	if removeName != "" {
		updatef = func(rotation *api.Rotation) error {
			return rotation.DeleteNeed(removeName)
		}
	} else {
		if level == "" || skill == "" || min == 0 {
			return "", errors.Errorf("--need requires skill, level, and min to be specified. %s",
				commandUsage("rotation need <rotation-name>", fs))
		}
		l := 0
		l, err = api.ParseLevel(level)
		if err != nil {
			return "", err
		}
		updatef = func(rotation *api.Rotation) error {
			rotation.ChangeNeed(name, store.Need{
				Skill: skill,
				Level: l,
				Min:   min,
				Max:   max,
			})
			return nil
		}
	}

	err = c.API.UpdateRotation(rotation, updatef)
	if err != nil {
		return "", err
	}

	return "Updated rotation " + api.MarkdownRotationWithDetails(rotation), nil
}

func withRotationCreate(fs *pflag.FlagSet, start *string, period *api.Period) {
	fs.StringVar(start, "start", "", fmt.Sprintf("rotation start date formatted as %s. It must be provided at creation and **can not be modified** later.", api.DateFormat))
	fs.Var(period, "period", "rotation period 1w, 2w, or 1m")
}

func withRotationOptions(fs *pflag.FlagSet, size *int, paddingWeeks *int) {
	fs.IntVar(size, "size", 0, "target number of people in each shift. 0 (default) means unlimited, based on needs")
	fs.IntVar(paddingWeeks, "padding", 0, "makes each user's shift  padded by this many weeks of unavailability, on each side")
}

func withRotationID(fs *pflag.FlagSet, rotationID, rotationName *string) {
	fs.StringVar(rotationID, "rotation-id", "", "specify rotation ID")
	fs.StringVar(rotationName, "rotation-name", "", "specify rotation name")
}

func (c *Command) parseRotationFlagsAndLoad(fs *pflag.FlagSet, parameters []string, subcommand string) (*api.Rotation, error) {
	var id, name string
	withRotationID(fs, &id, &name)
	err := fs.Parse(parameters)
	arg := fs.Arg(0)
	switch {
	case err != nil:
		return nil, errors.Errorf("**%s**\n\n%s", err.Error(), commandUsage(subcommand, fs))

	case id == "" && arg == "" && name == "":
		return nil, errors.Errorf("**Rotation is not specified**\n\n%s", commandUsage(subcommand, fs))

	case id != "" && arg != "" && name != "",
		id != "" && name != "",
		id != "" && arg != "",
		arg != "" && name != "":
		return nil, errors.Errorf("**rotation is specified multiple times**\n\n%s", commandUsage(subcommand, fs))

	case id != "":
		return c.API.LoadRotation(id)

	case arg != "", name != "":
		if name == "" {
			name = arg
		}
		return c.API.LoadRotationNamed(name)
	}
	return nil, errors.New("unreachable go silly")
}

func withRotationNeed(fs *pflag.FlagSet, name, skill, level *string, min, max *int, removeName *string) {
	fs.StringVar(name, "need", "", "update rotation need")
	fs.StringVar(skill, "skill", "", "if used with --need, indicates the needed skill")
	fs.StringVar(level, "level", "", "if used with --need, indicates the needed skill level")
	fs.IntVar(min, "min", 0, "if used with --need, indicates the minimum needed headcount")
	fs.IntVar(max, "max", 0, "if used with --need, indicates the maximum needed headcount")
	fs.StringVar(removeName, "remove-need", "", "remove a need from rotation")
}
