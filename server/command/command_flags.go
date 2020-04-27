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

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const intNoValue = int64(0xBAADBEEF)

func (c *Command) flagUsage() md.MD {
	usage := c.actualTrigger
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
	usage := fmt.Sprintf("`%s %s`", c.actualTrigger, strings.Join(subs, "|"))
	return md.Markdownf("Usage: %s\nUse `%s <subcommand> --help` for more info.",
		usage, c.actualTrigger)
}

func (c *Command) parse(parameters []string) error {
	err := c.flags().Parse(parameters)
	if err != nil {
		return err
	}

	if (*c.now).IsZero() {
		now := types.NewTime(time.Now())
		c.now = &now
	}
	return nil
}

func (c *Command) flags() *pflag.FlagSet {
	if c.fs == nil {
		c.fs = pflag.NewFlagSet("", pflag.ContinueOnError)
		c.fs.BoolVar(&c.outputJson, "json", false, "output as JSON")
		c.now, _ = c.withTimeFlag("now", "specify the transaction time (default: now)")
	}
	return c.fs
}

func (c *Command) withTimeFlag(flag, desc string) (*types.Time, error) {
	actingUser, err := c.SL.ActingUser()
	if err != nil {
		return nil, err
	}
	t := actingUser.Time(types.Time{})
	c.flags().Var(&t, flag, desc)
	return &t, nil
}

func (c *Command) withFlagRotation() {
	c.flags().StringP("rotation", "r", "", "rotation reference")
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
		user, err := c.SL.LoadMattermostUserByUsername(arg)
		if err != nil {
			return nil, err
		}
		mattermostUserIDs.Set(user.MattermostUserID)
	}

	return mattermostUserIDs, nil
}

func (c *Command) resolveRotationUsernames() (types.ID, *types.IDSet, error) {
	ref, _ := c.flags().GetString("rotation")
	usernames := []string{}
	rotationID := types.ID(ref)

	for _, arg := range c.flags().Args() {
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
	args := c.flags().Args()
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
	ref, _ := c.flags().GetString("rotation")
	rotationID := types.ID(ref)
	if ref == "" {
		if len(c.flags().Args()) < 1 {
			return "", errors.New("no rotation specified")
		}
		rotationID, err = c.SL.ResolveRotationName(c.flags().Arg(0))
		if err != nil {
			return "", err
		}
	}
	return rotationID, nil
}
