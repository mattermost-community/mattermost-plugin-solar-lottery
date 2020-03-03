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
	- [x] archive
	- [x] debug-delete
	- [x] list
	- [x] new
	- [x] show

- [ ] task
	- [ ] param min --skill <s-l> (--count | --clear)
	- [ ] param max --skill <s-l> (--count | --clear)
	- [ ] param grace --duration 
	- [ ] param shift --starting --period
	- [ ] param ticket
	- [ ] new ticket
	- [ ] new shift
	- [ ] assign [--auto aka fill]
	- [ ] close
	- [ ] debug-delete
	- [ ] list --pending | --in-progress | --all
	- [ ] start
	- [ ] unassign

- [x] skill
	- [x] delete
	- [x] list
	- [x] new

- [x] user: manage users.
	- [x] join ROT [@user...] --starting
	- [x] leave ROT [@user...]
	- [x] disqualify --skill 
	- [x] qualify --skill 
	- [x] show 
	- [x] unavailable: --start --finish [--clear] 

	`
	return resp, nil
}
