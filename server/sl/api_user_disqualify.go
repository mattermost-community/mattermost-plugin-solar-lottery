// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InDisqualify struct {
	MattermostUserIDs *types.IDSet
	Skills            []string
}

func (sl *sl) Disqualify(params InDisqualify) (*OutQualify, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("Disqualify", params),
		withValidSkillNames(params.Skills...),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	err = sl.disqualify(users, params.Skills)
	if err != nil {
		return nil, err
	}

	out := &OutQualify{
		Users: users,
		MD:    md.Markdownf("removed skill(s) %s from %s.", params.Skills, users.Markdown()),
	}
	sl.logAPI(out)
	return out, nil
}
