// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
)

func (c *Command) info(parameters []string) (string, error) {
	resp := fmt.Sprintf("Mattermost Solar Lottery plugin version: %s, "+
		"[%s](https://github.com/mattermost/%s/commit/%s), built %s\n",
		c.Config.PluginVersion,
		c.Config.BuildHashShort,
		config.Repository,
		c.Config.BuildHash,
		c.Config.BuildDate)

	resp += `
- [x] info: display this.

- [x] rotation
	- [x] add
	- [x] archive
	- [x] debug-delete
	- [x] forecast
	- [x] guess
	- [x] join
	- [x] leave
	- [x] list
	- [x] need (add/delete)
	- [x] show
	- [x] update

- [ ] shift
	- [ ] add
	- [ ] commit: closes the shift for volunteers, notifies selected users.
	- [ ] debug-delete
	- [ ] fill: evaluates shift readiness, autofills.
	- [ ] finish: finishes a shift.
	- [ ] join: add user(s) to shift.
	- [ ] leave: remove user(s) from shift.
	- [ ] list
	- [ ] show
	- [ ] start: starts a shift.

- [x] skill
	- [x] add
	- [x] list
	- [x] delete

- [x] user: manage my profile.
	- [-TODO-] forecast
	- [x] show [--users] 
	- [x] unavailable: --from --to [--clear] [--type=unavailable]
	- [x] qualify --skill --level --users
	- [x] disqualify --skill --users
`
	return resp, nil
}
