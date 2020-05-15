// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InJoinRotation struct {
	MattermostUserIDs *types.IDSet
	RotationID        types.ID
	Starting          types.Time
}

type OutJoinRotation struct {
	md.MD
	Modified *Users
}

func (sl *sl) JoinRotation(params InJoinRotation) (*OutJoinRotation, error) {
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
		modified, err = sl.joinRotation(users, r, params.Starting)
		return err
	})
	if err != nil {
		return nil, err
	}

	out := &OutJoinRotation{
		Modified: modified,
		MD:       md.Markdownf("added %s to %s.", modified.MarkdownWithSkills(), r.Markdown()),
	}
	sl.logAPI(out)
	return out, nil
}
