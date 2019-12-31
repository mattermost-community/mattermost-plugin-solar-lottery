// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
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
	- [x] debug-rotation
	- [-] debug-shift
	- [x] skill

- [ ] show
	- [x] rotation
	- [ ] shift
	- [x] user

- [ ] list
	- [x] rotation
	- [-] shift
	- [x] skill
	- [ ] user

- [x] rotation
	- [x] archive: archive rotation.
	- [x] debug-delete
	- [x] show: same as "show rotation"
	- [x] add: same as "add rotation"
	- [x] update

- [ ] forecast
	- [x] guess
	- [x] rotation
	- [ ] user
	
- [ ] shift
	- [-] join: add user(s) to shift.
	- [ ] leave: remove user(s) from shift.
	- [ ] fill: evaluates shift readiness, autofills.
	- [-] commit: closes the shift for volunteers, notifies selected users.
	- [-] finish: finishes a shift.
	- [-] start: starts a shift.
	- [x] add

- [x] user: manage my profile.
	- [x] show [--users] 
	- [x] unavailable: --from --to [--clear] [--type=unavailable]
	- [x] qualify --skill --level --users
	- [x] disqualify --skill --users

- [x] need (add/delete)

- [x] join: add user(s) to rotation.

- [x] leave: remove user(s) from rotation.

- [x] info: display plugin information.
	
	`
	return resp, nil
}
