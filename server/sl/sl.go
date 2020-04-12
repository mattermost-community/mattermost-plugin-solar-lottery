// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

const (
	ctxActingUserID   = "ActingUserID"
	ctxActingUsername = "ActingUsername"
	ctxAPI            = "API"
	ctxInput          = "Input"
	ctxFill           = "Fill"
	ctxForce          = "Force"
	ctxInterval       = "Interval"
	ctxRotationID     = "RotationID"
	ctxSkill          = "Skill"
	ctxSkillLevel     = "SkillLevel"
	ctxSourceName     = "SourceName"
	ctxStarting       = "Starting"
	ctxTaskID         = "TaskID"
	ctxUnavailable    = "Unavailable"
	ctxUserIDs        = "UserIDs"
	ctxUsernames      = "Usernames"
	ctxUsers          = "Users"
)

var ErrMultipleResults = errors.New("multiple results found")
var ErrAlreadyExists = errors.New("already exists")

type TaskService interface {
	AssignTask(InAssignTask) (*OutAssignTask, error)
	UnassignTask(InAssignTask) (*OutAssignTask, error)
	FillTask(InAssignTask) (*OutAssignTask, error)
	LoadTask(types.ID) (*Task, error)
	TransitionTask(params InTransitionTask) (*OutTransitionTask, error)
	CreateTicket(InCreateTicket) (*OutCreateTask, error)
	CreateShift(InCreateShift) (*OutCreateTask, error)
}

type SkillService interface {
	ListKnownSkills() (*types.IDSet, error)
	AddKnownSkill(types.ID) error
	DeleteKnownSkill(types.ID) error
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

type AutopilotService interface {
	RunAutopilot(in *InRunAutopilot) (*OutRunAutopilot, error)
}

type SL interface {
	RotationService
	SkillService
	UserService
	TaskService
	AutopilotService

	PluginAPI
	bot.Logger

	ActingUser() (*User, error)
	Config() *config.Config

	LoadUsers(mattermostUserIDs *types.IDSet) (*Users, error)
	LoadMattermostUserByUsername(username string) (*User, error)
}

type sl struct {
	*Service
	bot.Logger

	conf *config.Config

	// set by Service.ActingAs.
	actingMattermostUserID types.ID

	// set by withActingUser or withActingUserExpanded.
	actingUser *User

	// Stack of loggers
	loggers []bot.Logger
}

func (sl *sl) Config() *config.Config {
	return sl.conf
}

func (sl *sl) logAPI(msg md.Markdowner) {
	sl.Infof("%s: %s", sl.actingUser.Markdown(), msg.Markdown())
}
