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
	Type             types.ID
	Min              Needs         `json:`
	Max              Needs         `json:`
	Grace            time.Duration `json:"grace,omitempty"`
	TicketSeq        int           `json:",omitempty"`
	ShiftDescription string        `json:",omitempty"`
	ShiftStart       types.Time    `json:",omitempty"`
	ShiftPeriod      types.Period  `json:",omitempty"`
}

func NewTaskMaker() *TaskMaker {
	return &TaskMaker{
		Min: NewNeeds(),
		Max: NewNeeds(),
	}
}

func (maker *TaskMaker) NewTicket(rotationID types.ID, defaultID string) *Task {
	maker.TicketSeq++
	t := NewTask(rotationID)
	if defaultID == "" {
		defaultID = strconv.Itoa(maker.TicketSeq)
	}
	t.TaskID = types.ID(string(rotationID) + "-" + defaultID)
	// TODO do I need to clone these?
	t.Min = maker.Min
	t.Max = maker.Max
	t.Grace = maker.Grace
	return t
}

// func (maker *TaskMaker) NewShift(prefix string, t types.Time) *Task {
// 	task := maker.NewTicket(prefix, "")
// 	return task
// }

func (maker TaskMaker) MarkdownBullets() string {
	out := fmt.Sprintf("  - Type: **%s**\n", maker.Type)
	out += fmt.Sprintf("  - Min: **%s**\n", maker.Min.Markdown())
	out += fmt.Sprintf("  - Max: **%v**\n", maker.Max.Markdown())
	out += fmt.Sprintf("  - Grace: %v\n", maker.Grace)
	if maker.Type == ShiftMaker {
		out += fmt.Sprintf("  - Shift period: **%v**\n", maker.ShiftPeriod.String())
		out += fmt.Sprintf("  - Started: **%v**\n", maker.ShiftStart)
	}
	return out
}
