// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) taskTransition(to types.ID) func([]string) (md.MD, error) {
	return func(parameters []string) (md.MD, error) {
		c.assureFS()
		err := c.fs.Parse(parameters)
		if err != nil {
			return c.flagUsage(), err
		}
		taskID, _, err := c.resolveTaskIDUsernames()
		if err != nil {
			return "", err
		}

		return c.normalOut(c.SL.TransitionTask(sl.InTransitionTask{
			TaskID: taskID,
			State:  to,
			Time:   types.NewTime(),
		}))
	}
}
