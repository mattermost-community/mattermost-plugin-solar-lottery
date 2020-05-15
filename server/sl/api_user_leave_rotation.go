// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (sl *sl) LeaveRotation(params InJoinRotation) (*OutJoinRotation, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("LeaveRotation", params),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	modified := NewUsers()
	r, err := sl.UpdateRotation(params.RotationID, func(r *Rotation) error {
		modified, err = sl.leaveRotation(users, r)
		return err
	})
	if err != nil {
		return nil, err
	}

	out := &OutJoinRotation{
		Modified: modified,
		MD:       md.Markdownf("removed %s from %s.", modified.MarkdownWithSkills(), r.Markdown()),
	}
	sl.logAPI(out)
	return out, nil
}
