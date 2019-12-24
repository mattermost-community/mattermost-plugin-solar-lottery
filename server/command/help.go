// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/spf13/pflag"
)

func (c *Command) help(parameters []string) (string, error) {
	resp := fmt.Sprintf("Mattermost Solar Lottery plugin version: %s, "+
		"[%s](https://github.com/mattermost/%s/commit/%s), built %s\n",
		c.Config.PluginVersion,
		c.Config.BuildHashShort,
		config.Repository,
		c.Config.BuildHash,
		c.Config.BuildDate)

	resp += `
- [-] add
	- [x] rotation
	- [-] shift
	- [x] skill

- [ ] delete
	- [-] debug-rotation
	- [-] debug-shift
	- [-] skill

- [ ] show
	- [x] rotation
	- [ ] shift
	- [x] user

- [ ] forecast
	- [ ] schedule
	- [ ] heat

- [ ] list
	- [x] rotation
	- [-] shift
	- [x] skill
	- [ ] user

- [ ] rotation
	- [x] archive: archive rotation.
	- [ ] debug-delete
	- [ ] show: same as "show rotation"
	- [ ] add: same as "add rotation"
	- [x] update
	
- [ ] need

- [ ] join: add user(s) to rotation.

- [ ] leave: remove user(s) from rotation.

- [ ] shift
	- [-] join: add user(s) to shift.
	- [ ] leave: remove user(s) from shift.
	- [ ] fill: evaluates shift readiness, autofills.
	- [-] commit: closes the shift for volunteers, notifies selected users.
	- [-] finish: finishes a shift.
	- [-] start: starts a shift.
	- [x] add

- [x] info: display plugin information.

- [ ] user: manage my profile.
	- [x] show [--users] 
	- [ ] unavailable: --from --to [--clear] [--type=unavailable]
	- [-] qualify --skill --level --users
	- [ ] disqualify --skill --users
`
	return resp, nil
}

func subusage(command string, fs *pflag.FlagSet) string {
	if fs == nil {
		return fmt.Sprintf("Usage:\n```\n/%s %s```\n", config.CommandTrigger, command)
	}
	return fmt.Sprintf("Usage:\n```\n/%s %s [flags...]\n\n%s```\n",
		config.CommandTrigger, command, fs.FlagUsages())
}
