// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

// lotto rotation issue-source delete ROT KEY
// lotto rotation issue-source limit ROT KEY --clear --skill webapp-3
// lotto rotation issue-source limit ROT KEY --limit --count 1 --skill webapp-3
// lotto rotation issue-source put ROT KEY --grace 720h
// lotto rotation issue-source require ROT KEY --clear --skill webapp[-3]
// lotto rotation issue-source require ROT KEY --count 2 --skill webapp-2
func (c *Command) rotationIssueSource(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		commandDelete:  c.rotationIssueSourceDelete,
		commandRequire: c.rotationIssueSourceRequire,
		// commandLimit:   c.rotationIssueSourceLimit,
		// commandPut:     c.rotationIssueSourcePut,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) rotationIssueSourceDelete(parameters []string) (string, error) {
	fs := newFS()
	fRotation(fs)
	jsonOut := fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}
	r, err := c.loadRotation(fs)
	if err != nil {
		return "", err
	}

	err = c.SL.DeleteIssueSource(r, types.ID(fs.Arg(1)))
	if err != nil {
		return "", err
	}

	if *jsonOut {
		return md.JSONBlock(r), nil
	}
	return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), r.Markdown()), nil
}

func (c *Command) rotationIssueSourceRequire(parameters []string) (string, error) {
	// fs := newFS()
	// fRotation(fs)
	// jsonOut := fJSON(fs)
	// clear := fClear(fs)
	// count := fCount(fs)
	// skillLevel := fSkillLevel(fs)
	// err := fs.Parse(parameters)
	// if err != nil {
	// 	return c.flagUsage(fs), err
	// }
	// r, err := c.loadRotation(fs)
	// if err != nil {
	// 	return "", err
	// }

	// r.
	// 	err = c.SL.DeleteIssueSource(r, fs.Arg(1))
	// if err != nil {
	// 	return "", err
	// }

	// if *jsonOut {
	// 	return md.JSONBlock(r), nil
	// }
	// return fmt.Sprintf("%s deleted from rotation %s", fs.Arg(1), r.Markdown()), nil
	return "<><>", nil
}
