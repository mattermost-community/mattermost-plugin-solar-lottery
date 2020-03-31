// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

var ErrMultipleResults = errors.New("multiple results found")
var ErrAlreadyExists = errors.New("already exists")

type TaskService interface {
	AssignTask(InAssignTask) (*OutAssignTask, error)
	UnassignTask(InAssignTask) (*OutAssignTask, error)
	FillTask(InAssignTask) (*OutAssignTask, error)
	LoadTask(types.ID) (*Task, error)
	TransitionTask(params InTransitionTask) (*OutTransitionTask, error)
	MakeTicket(InMakeTicket) (*OutMakeTask, error)
	MakeShift(InMakeShift) (*OutMakeTask, error)
}

type UserService interface {
	AddToCalendar(InAddToCalendar) (*OutCalendar, error)
	ClearCalendar(InClearCalendar) (*OutCalendar, error)
	Disqualify(InQualify) (*OutQualify, error)
	JoinRotation(InJoinRotation) (*OutJoinRotation, error)
	LeaveRotation(InJoinRotation) (*OutJoinRotation, error)
	Qualify(InQualify) (*OutQualify, error)
}

type RotationService interface {
	AddRotation(*Rotation) error
	ArchiveRotation(rotationID types.ID) (*Rotation, error)
	DebugDeleteRotation(rotationID types.ID) error
	LoadActiveRotations() (*types.IDSet, error)
	LoadRotation(rotationID types.ID) (*Rotation, error)
	MakeRotation(rotationName string) (*Rotation, error)
	ResolveRotationName(string) (types.ID, error)
	UpdateRotation(rotationID types.ID, updatef func(*Rotation) error) (*Rotation, error)
}
