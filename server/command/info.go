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
	- [x] join
	- [x] leave
	- [x] list
	- [x] new
	- [x] show

- [x] skill
	- [x] delete
	- [x] list
	- [x] new

- [x] user: manage users.
	- [x] disqualify --skill 
	- [x] qualify --skill --level 
	- [x] show: show users
	- [x] unavailable: --start --finish [--clear] 

- [ ] task
	- [ ] assign [--auto aka fill]
	- [ ] close
	- [ ] debug-delete
	- [ ] list (by rotation, by task queue)
	- [ ] new issue
	- [ ] new shift 
	- [ ] show
	- [ ] start
	- [ ] unassign

	`
	return resp, nil
}
