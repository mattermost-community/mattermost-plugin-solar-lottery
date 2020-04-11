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

- [ ] autopilot

- [ ] task
	- [ ] debug-delete
	- [ ] list --pending | --in-progress | --all
	- [x] assign
	- [x] close
	- [x] fill
	- [x] new shift
	- [x] new ticket
	- [x] schedule
	- [x] show ROT#id
	- [x] start
	- [x] unassign

- [x] rotation
	- [x] archive ROT
	- [x] debug-delete ROT
	- [x] list
	- [x] new ROT
	- [x] show ROT
	- [x] param grace --duration 
	- [x] param max --skill <s-l> (--count | --clear)
	- [x] param min --skill <s-l> (--count | --clear)
	- [x] param shift --starting --period
	- [x] param ticket
	
- [x] skill
	- [x] delete SKILL
	- [x] list
	- [x] new SKILL

- [x] user: manage users.
	- [x] disqualify [@user...] --skill 
	- [x] join ROT [@user...] --starting
	- [x] leave ROT [@user...]
	- [x] qualify [@user...] --skill 
	- [x] show [@user...]
	- [x] unavailable: [@user...] --start --finish [--clear] 


	`
	return resp, nil
}
