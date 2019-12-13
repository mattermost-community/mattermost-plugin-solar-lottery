// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func (c *Command) join(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]

	asUsername := ""
	graceShifts := 0
	s := flag.NewFlagSet("join", flag.ContinueOnError)
	s.StringVar(&asUsername, "user", "", "add another user to the rotation.")
	s.IntVar(&graceShifts, "grace", 0, "start with N grace shifts.")
	err := s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}

	err = c.API.JoinRotation(rotationName, graceShifts, asUsername)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Joined rotation %s", rotationName), nil
}

func (c *Command) leave(parameters ...string) (string, error) {
	if len(parameters) == 0 {
		return "", errors.New("invalid syntax TODO")
	}
	rotationName := parameters[0]

	asUsername := ""
	s := flag.NewFlagSet("leave", flag.ContinueOnError)
	s.StringVar(&asUsername, "user", "", "add another user to the rotation.")
	err := s.Parse(parameters[1:])
	if err != nil {
		return "", err
	}

	err = c.API.LeaveRotation(rotationName, asUsername)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Left rotation %s", rotationName), nil
}
