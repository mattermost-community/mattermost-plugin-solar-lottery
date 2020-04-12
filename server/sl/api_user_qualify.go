// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InQualify struct {
	MattermostUserIDs *types.IDSet
	SkillLevel        SkillLevel
}

type OutQualify struct {
	md.MD
	Users *Users
}

func (sl *sl) Qualify(params InQualify) (*OutQualify, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("Qualify", params),
		// NOT restricted to: withValidSkillName(&params.SkillLevel.Skill),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	err = sl.qualify(users, params.SkillLevel)
	if err != nil {
		return nil, err
	}

	out := &OutQualify{
		Users: users,
		MD:    md.Markdownf("added skill %s to %s.", params.SkillLevel, users.Markdown()),
	}
	sl.logAPI(out)
	return out, nil
}
