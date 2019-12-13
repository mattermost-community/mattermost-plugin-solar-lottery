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
- [x] join: add a user to rotation.
	- [x] --user @user: add someone else to rotation
- [x] leave: remove a user from rotation.
	- [x] --user @user: remove someone else from rotation
- [ ] rotation: manage rotations.
	- [x] list: list rotations.
	- [x] show: show rotation's details.
		- [x] --schedule: lists currently scheduled shifts.
		- [x] --auto: adds auto-filled shifts to --schedule.
		- [x] --shifts: list this many shifts forward in time
	- [x] add: add a new rotation.
	- [x] delete: delete a rotation.
	- [ ] update: modiy rotation's settings.
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
- [x] me: manage my profile
	- [x] show: display the profile.
	- [x] skills: manage skill levels.
`
	return resp, nil
}
