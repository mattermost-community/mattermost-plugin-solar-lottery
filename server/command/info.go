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

- [ ] autopilot [--now=datetime]

- [ ] task
	- [ ] debug-delete
	- [ ] list --pending | --scheduled | --started | --finished
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
	- [x] set autopilot ROT 
		- [x] --off
		- [x] --create --create-prior[=28d]
		- [x] --schedule --schedule-prior[=7d]
		- [x] --start-finish
		- [x] --notify-start-prior[=3d]
		- [x] --notify-finish-prior[=3d]
		- [x] --run=time
	- [x] set fill
		- [ ] --beginning
		- [ ] --period
		- [ ] --seed
		- [ ] --fuzz
	- [x] set limit --skill <s-l> (--count | --clear)
	- [x] set require --skill <s-l> (--count | --clear)
	- [x] set task
		- [ ] --type=(shift|ticket)
		- [ ] --duration
		- [ ] --grace
	
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
