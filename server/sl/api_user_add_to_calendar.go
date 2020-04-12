// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InAddToCalendar struct {
	MattermostUserIDs *types.IDSet
	Unavailable       *Unavailable
}

type OutCalendar struct {
	Users *Users
	md.MD
}

func (sl *sl) AddToCalendar(params InAddToCalendar) (*OutCalendar, error) {
	users := NewUsers()
	err := sl.Setup(
		pushAPILogger("AddToCalendar", params),
		withExpandedUsers(&params.MattermostUserIDs, users),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	for _, user := range users.AsArray() {
		sl.addUserUnavailable(user, params.Unavailable)
	}

	out := &OutCalendar{
		Users: users,
		MD: md.Markdownf("added unavailable event %s to %s",
			sl.actingUser.MarkdownUnavailable(params.Unavailable), users.Markdown()),
	}
	sl.logAPI(out)
	return out, nil
}
