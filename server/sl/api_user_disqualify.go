// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (sl *sl) Disqualify(params InQualify) (*OutQualify, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("Disqualify", params),
		withValidSkillName(&params.SkillLevel.Skill),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	err = sl.disqualify(users, params.SkillLevel.Skill)
	if err != nil {
		return nil, err
	}

	out := &OutQualify{
		Users: users,
		MD:    md.Markdownf("removed skill %s from %s.", params.SkillLevel.Skill, users.Markdown()),
	}
	sl.logAPI(out)
	return out, nil
}
