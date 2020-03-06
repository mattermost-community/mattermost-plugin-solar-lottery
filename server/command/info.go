// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) info(parameters []string) (md.MD, error) {
	resp := md.Markdownf("Mattermost Solar Lottery plugin version: %s, "+
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
	- [x] param min --skill <s-l> (--count | --clear)
	- [x] param max --skill <s-l> (--count | --clear)
	- [x] param grace --duration 
	- [x] param shift --starting --period
	- [x] param ticket
	- [x] new ticket
	- [ ] new shift
	- [x] assign
	- [ ] assign --fill
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
