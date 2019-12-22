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
- [ ] shift:
	- [x] commit: closes the shift for volunteers, notifies selected users.
	- [ ] fill: adds users to a shift.
	- [x] finish: finishes a shift.
	- [x] list: list shifts for a rotation.
	- [x] open: creates a shift and invites users to volunteer.
	- [x] start: starts a shift.
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
- [x] skill: manage known skills.
	- [x] list: list skills.
	- [x] add: add a new skill.
	- [x] delete: delete a skill.
- [x] user: manage my profile.
	- [x] show: show user details.
- [x] volunteer: add user(s) to an open shift
`
	return resp, nil
}
