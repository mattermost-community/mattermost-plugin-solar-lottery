// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InMakeTicket struct {
	RotationID  types.ID
	Summary     string
	Description string
}

type OutMakeTicket struct {
	md.MD
	Task *Task
}

func (sl *sl) MakeTicket(params InMakeTicket) (*OutMakeTicket, error) {
	err := sl.Setup(pushAPILogger("MakeTicket", params))
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	var t *Task
	_, err = sl.UpdateRotation(params.RotationID, func(r *Rotation) error {
		t = r.TaskMaker.newTicket(r, "")
		t.Summary = params.Summary
		t.Description = params.Description
		var id types.ID
		id, err = sl.Store.Entity(KeyTask).NewID(string(t.TaskID))
		if err != nil {
			return err
		}
		t.TaskID = id

		err = sl.storeTask(t)
		if err != nil {
			return err
		}
		r.TaskIDs.Set(t.TaskID)
		if r.pending != nil {
			r.pending.Set(t)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := &OutMakeTicket{
		MD:   md.Markdownf("created ticket %s.", t.Markdown()),
		Task: t,
	}
	sl.LogAPI(out)
	return out, nil
}
