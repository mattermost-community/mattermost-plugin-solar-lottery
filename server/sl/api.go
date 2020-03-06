// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

type TaskService interface {
	MakeTicket(InMakeTicket) (*OutMakeTicket, error)
	AssignTask(InAssignTask) (*OutAssignTask, error)
}

type UserService interface {
	AddToCalendar(InAddToCalendar) (*OutCalendar, error)
	ClearCalendar(InClearCalendar) (*OutCalendar, error)
	Disqualify(InDisqualify) (*OutQualify, error)
	Qualify(InQualify) (*OutQualify, error)
	JoinLeaveRotation(InJoinLeaveRotation) (*OutJoinLeaveRotation, error)
}
