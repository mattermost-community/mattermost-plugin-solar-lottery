// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const (
	TicketMaker = types.ID("TicketMaker")
	ShiftMaker  = types.ID("ShiftMaker")
)

type TaskMaker struct {
	Type                  types.ID
	Require               *Needs        `json:",omitempty"`
	Limit                 *Needs        `json:",omitempty"`
	Grace                 time.Duration `json:",omitempty"`
	TicketSeq             int           `json:",omitempty"`
	TicketDefaultDuration time.Duration `json:",omitempty"`
	ShiftDescription      string        `json:",omitempty"`
	ShiftStart            types.Time    `json:",omitempty"`
	ShiftPeriod           types.Period  `json:",omitempty"`
}

func NewTaskMaker() *TaskMaker {
	return &TaskMaker{
		Require: NewNeeds(),
		Limit:   NewNeeds(),
	}
}

func (maker *TaskMaker) newTicket(r *Rotation, defaultID string) *Task {
	maker.TicketSeq++
	t := NewTask(r.RotationID)
	if defaultID == "" {
		defaultID = strconv.Itoa(maker.TicketSeq)
	}
	t.TaskID = types.ID(r.Name() + "#" + defaultID)
	// TODO do I need to clone these?
	t.Require = maker.Require
	t.Limit = maker.Limit
	t.Grace = maker.Grace
	return t
}

func (maker TaskMaker) MarkdownBullets() string {
	out := fmt.Sprintf("  - Type: **%s**\n", maker.Type)
	out += fmt.Sprintf("  - Require: **%s**\n", maker.Require.Markdown())
	out += fmt.Sprintf("  - Limit: **%v**\n", maker.Limit.Markdown())
	out += fmt.Sprintf("  - Grace: %v\n", maker.Grace)
	if maker.Type == ShiftMaker {
		out += fmt.Sprintf("  - Shift period: **%v**\n", maker.ShiftPeriod.String())
		out += fmt.Sprintf("  - Started: **%v**\n", maker.ShiftStart)
	}
	return out
}
