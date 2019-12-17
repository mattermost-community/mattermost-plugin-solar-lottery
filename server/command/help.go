// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
)

func (c *Command) help(parameters ...string) (string, error) {
	resp := fmt.Sprintf("Mattermost Solar Lottery plugin version: %s, "+
		"[%s](https://github.com/mattermost/%s/commit/%s), built %s\n",
		c.Config.PluginVersion,
		c.Config.BuildHashShort,
		config.Repository,
		c.Config.BuildHash,
		c.Config.BuildDate)
	resp += `
- [x] info: display plugin information.
- [x] join: add user(s) to rotation.
- [x] leave: remove user(s) from rotation.
- [x] rotation: manage rotations.
	- [x] list: list rotations.
	- [x] show: show rotation details.
	- [x] forecast: show rotation forecast.
	- [x] create: create a new rotation.
	- [x] archive: archives a rotation.
	- [x] update: modiy rotation's parameters.
- [ ] shift:
	- [ ] schedule: creates/updates a scheduled shift for a rotation.
		- [ ] --in: schedules a user into shift.
		- [ ] --out: schedules a user in  shift.
		- [ ] --create: create a new shift
	- [ ] delete: deletes (resets) a scheduled shift.
	- [ ] start: starts a shift (??).
	- [ ] finish: finishes a shift.
- [x] skill: manage known skills.
	- [x] list: list skills.
	- [x] add: add a new skill.
	- [x] delete: delete a skill.
- [x] user: manage my profile.
	- [x] show: show user details.
`
	return resp, nil
}
