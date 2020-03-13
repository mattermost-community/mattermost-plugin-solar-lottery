// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"

type TaskService interface {
	LoadTask(types.ID) (*Task, error)
	MakeTicket(InMakeTicket) (*OutMakeTicket, error)
	AssignTask(InAssignTask) (*OutAssignTask, error)
	FillTask(InAssignTask) (*OutAssignTask, error)
}

type UserService interface {
	AddToCalendar(InAddToCalendar) (*OutCalendar, error)
	ClearCalendar(InClearCalendar) (*OutCalendar, error)
	Disqualify(InDisqualify) (*OutQualify, error)
	Qualify(InQualify) (*OutQualify, error)
	JoinLeaveRotation(InJoinLeaveRotation) (*OutJoinLeaveRotation, error)
}
