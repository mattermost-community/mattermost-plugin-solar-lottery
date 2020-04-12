// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

type InClearCalendar struct {
	MattermostUserIDs *types.IDSet
	Interval          types.Interval
}

func (sl *sl) ClearCalendar(params InClearCalendar) (*OutCalendar, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("ClearCalendar", params),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	for _, user := range users.AsArray() {
		cleared := user.ClearUnavailable(params.Interval, "")
		if len(cleared) == 0 {
			continue
		}

		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to update user %s", user.Markdown())
		}
	}

	out := &OutCalendar{
		Users: users,
		MD:    md.Markdownf("deleted events %v from users %s.", params.Interval, users.MarkdownWithSkills()),
	}
	sl.logAPI(out)
	return out, nil
}
