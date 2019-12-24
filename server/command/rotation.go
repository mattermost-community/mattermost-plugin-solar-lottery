// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
)

func (c *Command) rotation(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandAdd:         c.addRotation,
		commandArchive:     c.archiveRotation,
		commandDebugDelete: c.debugDeleteRotation,
		commandShow:        c.showRotation,
		commandUpdate:      c.updateRotation,
	}

	return c.handleCommand(subcommands, parameters,
		"Usage: `rotation add|archive|show|update]`. Use `rotation subcommand --help` for more information.")
}

func (c *Command) addRotation(parameters []string) (string, error) {
	var rotationName, start string
	var period api.Period
	var size, paddingWeeks int
	fs := flag.NewFlagSet("addRotation", flag.ContinueOnError)
	withRotationCreateFlags(fs, &start, &period)
	withRotationUpdateFlags(fs, &size, &paddingWeeks)
	fs.StringVar(&rotationName, flagRotation, "", "specify rotation name")
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("add rotation", fs), err
	}
	if rotationName == "" {
		return subusage("add rotation", fs), errors.Errorf("must specify rotation name, use `--%s`", flagRotation)
	}

	rotation, err := c.API.MakeRotation(rotationName)
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

func (c *Command) archiveRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := flag.NewFlagSet("archiveRotation", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("rotation archive", fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	err = c.API.ArchiveRotation(rotation)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to archive %s", rotation.Name)
	}

	return "Deleted rotation " + rotation.Name, nil
}

func (c *Command) debugDeleteRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := flag.NewFlagSet("debugDeleteRotation", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("delete debug-rotation", fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}

	err = c.API.DebugDeleteRotation(rotationID)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to delete %s", rotationID)
	}

	return "Deleted rotation " + rotationID, nil
}

func (c *Command) showRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	fs := flag.NewFlagSet("showRotation", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("show rotation", fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	return api.MarkdownRotationWithDetails(rotation), nil
}

func (c *Command) listRotations(parameters []string) (string, error) {
	if len(parameters) > 0 {
		return subusage("list rotations", nil), errors.New("unexpected parameters")
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

func (c *Command) updateRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	var size, paddingWeeks int
	fs := flag.NewFlagSet("updateRotation", flag.ContinueOnError)
	withRotationFlags(fs, &rotationID, &rotationName)
	withRotationUpdateFlags(fs, &size, &paddingWeeks)
	err := fs.Parse(parameters)
	if err != nil {
		return subusage("update rotation", fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
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

func withRotationCreateFlags(fs *pflag.FlagSet, start *string, period *api.Period) {
	fs.StringVar(start, flagStart, "", fmt.Sprintf("rotation start date formatted as %s. It must be provided at creation and **can not be modified** later.", api.DateFormat))
	fs.Var(period, flagPeriod, "rotation period 1w, 2w, or 1m")
}

func withRotationUpdateFlags(fs *pflag.FlagSet, size *int, paddingWeeks *int) {
	fs.IntVar(size, flagSize, 0, "target number of people in each shift. 0 (default) means unlimited, based on needs")
	fs.IntVar(paddingWeeks, flagPadding, 0, "makes each user's shift  padded by this many weeks of unavailability, on each side")
}

func withRotationFlags(fs *pflag.FlagSet, rotationID, rotationName *string) {
	fs.StringVar(rotationID, flagRotationID, "", "specify rotation ID")
	fs.StringVar(rotationName, flagRotation, "", "specify rotation name")
}

func (c *Command) parseRotationFlags(id, name string) (rotationID string, err error) {
	switch {
	case id == "" && name == "":
		return "", errors.New("rotation is not specified")

	case id != "" && name != "":
		return "", errors.New("rotation is specified multiple times")

	case id != "":
		return id, nil

	}
	//  name != "":
	rotationIDs, err := c.API.ResolveRotationName(name)
	if err != nil {
		return "", err
	}
	if len(rotationIDs) != 1 {
		return "", errors.Errorf("name %s is ambigous, please use --%s with one of %s", name, flagRotation, rotationIDs)
	}
	return rotationIDs[0], nil
}
