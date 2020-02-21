// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
)

func (c *Command) info(parameters []string) (string, error) {
	resp := fmt.Sprintf("Mattermost Solar Lottery plugin version: %s, "+
		"[%s](https://github.com/mattermost/%s/commit/%s), built %s\n",
		c.SL.Config().PluginVersion,
		c.SL.Config().BuildHashShort,
		constants.Repository,
		c.SL.Config().BuildHash,
		c.SL.Config().BuildDate)

	resp += `
- [x] info: display this.

- [x] rotation
	- [x] new
	- [x] archive
	- [ ] autopilot
	- [x] debug-delete
	- [x] forecast
	- [x] guess
	- [x] join
	- [x] leave
	- [x] list
	- [x] need (add/delete)
	- [x] show

- [ ] shift
	- [x] new
	- [x] debug-delete
	- [x] fill: evaluates shift readiness, autofills.
	- [x] finish: finishes a shift.
	- [x] join: add user(s) to shift.
	- [ ] leave: remove user(s) from shift.
	- [x] list
	- [ ] show
	- [x] start: starts a shift.

- [x] skill
	- [x] new
	- [x] list
	- [x] delete

- [x] user: manage my profile.
	- [x] forecast
	- [x] show [--users] 
	- [x] unavailable: --from --to [--clear] [--type=unavailable]
	- [x] qualify --skill --level --users
	- [x] disqualify --skill --users
`
	return resp, nil
}
