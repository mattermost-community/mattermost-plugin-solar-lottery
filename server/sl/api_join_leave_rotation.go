// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InJoinLeaveRotation struct {
	MattermostUserIDs *types.IDSet
	RotationID        types.ID
	Starting          types.Time
	Leave             bool
}

type OutJoinLeaveRotation struct {
	md.MD
	Modified *Users
}

func (sl *sl) JoinLeaveRotation(params InJoinLeaveRotation) (*OutJoinLeaveRotation, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("JoinRotation", params),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	modified := NewUsers()
	r, err := sl.UpdateRotation(params.RotationID, func(r *Rotation) error {
		if params.Leave {
			modified, err = sl.leaveRotation(users, r)
		} else {
			modified, err = sl.joinRotation(users, r, params.Starting)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	verb, prep := "added", "to"
	if params.Leave {
		verb, prep = "removed", "from"
	}
	out := &OutJoinLeaveRotation{
		Modified: modified,
		MD:       md.Markdownf("%s %s %s %s.", verb, modified.MarkdownWithSkills(), prep, r.Markdown()),
	}
	sl.LogAPI(out)
	return out, nil
}
