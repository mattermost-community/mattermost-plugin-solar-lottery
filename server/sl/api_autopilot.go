// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type InRunAutopilot struct {
	RotationID types.ID
	Time       types.Time
}

type OutRunAutopilot struct {
	md.MD
	Rotation *Rotation

	messages []md.MD
}

func (s *sl) RunAutopilot(in *InRunAutopilot) (*OutRunAutopilot, error) {
	r := NewRotation()
	out := &OutRunAutopilot{}

	autopilotOp := func(op func(*Rotation, types.Time) (md.Markdowner, error)) func(*sl) error {
		return func(*sl) error {
			msg, err := op(r, in.Time)
			if err != nil {
				return err
			}
			out.messages = append(out.messages, msg.Markdown())
			return nil
		}
	}

	err := s.Setup(
		pushAPILogger("RunAutopilot", in),
		withExpandedRotation(&in.RotationID, r),
		autopilotOp(s.autopilotRemindFinish),
		autopilotOp(s.autopilotFinish),
		autopilotOp(s.autopilotCreate),
		autopilotOp(s.autopilotFillSchedule),
		autopilotOp(s.autopilotRemindStart),
		autopilotOp(s.autopilotStart),
	)
	if err != nil {
		return nil, err
	}
	defer s.popLogger()

	out.Rotation = r
	out.MD = md.Markdownf("%s ran autopilot on %s for %v.", s.actingUser.Markdown(), r.Markdown(), in.Time)
	for _, msg := range out.messages {
		if msg != "" {
			out.MD += "\n  - " + msg
		}
	}

	s.logAPI(out)
	return out, nil
}
