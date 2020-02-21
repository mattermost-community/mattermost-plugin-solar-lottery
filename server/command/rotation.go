// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func (c *Command) rotation(parameters []string) (string, error) {
	subcommands := map[string]func([]string) (string, error){
		// commandAutopilot:   c.autopilotRotation,
		commandNew:         c.newRotation,
		commandArchive:     c.archiveRotation,
		commandDebugDelete: c.debugDeleteRotation,
		// commandForecast:    c.forecastRotation,
		// commandGuess:       c.guessRotation,
		commandJoin:  c.joinRotation,
		commandLeave: c.leaveRotation,
		commandList:  c.listRotations,
		// commandNeed:  c.rotationNeed,
		commandShow: c.showRotation,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) rotationUsers(fs *pflag.FlagSet) (*sl.Rotation, sl.UserMap, error) {
	ref, _ := fs.GetString(flagRotation)
	usernames := types.NewSet()
	rid := ref

	for _, arg := range fs.Args() {
		if strings.HasPrefix(arg, "@") {
			usernames.Add(arg)
		} else {
			if rid != "" {
				return nil, nil, errors.Errorf("rotation %s is already specified, cant't interpret %s", rid, arg)
			}
			rid = arg
		}
	}

	var err error
	var r *sl.Rotation
	if rid != "" {
		// explicit ref is used as is
		if ref == "" {
			rid, err = c.SL.ResolveRotation(rid)
			if err != nil {
				return nil, nil, err
			}
		}

		r, err = c.SL.LoadRotation(rid)
		if err != nil {
			return nil, nil, err
		}
	}

	users, err := c.users(usernames.AsArray())
	if err != nil {
		return nil, nil, err
	}

	return r, users, nil
}
