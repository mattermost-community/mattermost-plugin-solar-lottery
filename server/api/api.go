// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type API interface {
	bot.Logger
	PluginAPI

	Expander
	Forecaster
	Autopilot

	Rotations
	Shifts
	Skills
	Users
}

type Rotations interface {
	AddRotation(*Rotation) error
	ArchiveRotation(*Rotation) error
	DebugDeleteRotation(string) error
	LoadKnownRotations() (store.IDMap, error)
	LoadRotation(string) (*Rotation, error)
	MakeRotation(rotationName string) (*Rotation, error)
	ResolveRotationName(namePattern string) ([]string, error)
	UpdateRotation(*Rotation, func(*Rotation) error) error
}

type Autopilot interface {
	AutopilotRotation(rotation *Rotation, now time.Time) error
}

type Shifts interface {
	ListShifts(*Rotation, int, int) ([]*Shift, error)
	OpenShift(*Rotation, int) (*Shift, error)
	StartShift(*Rotation, int) (*Shift, error)
	FinishShift(*Rotation, int) (*Shift, error)
	DebugDeleteShift(*Rotation, int) error
	FillShift(*Rotation, int) (*Shift, UserMap, error)
	IsShiftReady(rotation *Rotation, shiftNumber int) (shift *Shift, ready bool, whyNot string, err error)
}

type Users interface {
	GetActingUser() (*User, error)
	LoadMattermostUsers(mattermostUsernames string) (UserMap, error)
	LoadStoredUsers(mattermostUserIDs store.IDMap) (UserMap, error)

	AddEvent(mattermostUsernames string, event Event) error
	DeleteEvents(mattermostUsernames string, startDate, endDate string) error
	Disqualify(mattermostUsernames, skillName string) error
	JoinRotation(mattermostUsernames string, rotation *Rotation, starting time.Time) (added UserMap, err error)
	JoinShift(mattermostUsernames string, rotation *Rotation, shiftNumber int) (*Shift, UserMap, error)
	LeaveRotation(mattermostUsernames string, rotation *Rotation) (deleted UserMap, err error)
	Qualify(mattermostUsernames, skillName string, level Level) error
}

type Forecaster interface {
	Guess(rotation *Rotation, startingShiftNumber, numShifts int) ([]*Shift, error)
	ForecastRotation(rotation *Rotation, startingShiftNumber, numShifts, sampleSize int) (*Forecast, error)
	ForecastUser(mattermostUsername string, rotation *Rotation, numShifts, sampleSize int, now time.Time) ([]float64, error)
}

type Skills interface {
	ListSkills() (store.IDMap, error)
	AddSkill(string) error
	DeleteSkill(string) error
}

type Expander interface {
	ExpandUserMap(UserMap) error
	ExpandUser(*User) error
	ExpandRotation(*Rotation) error
}

type PluginAPI interface {
	GetMattermostUser(mattermostUserID string) (*model.User, error)
	GetMattermostUserByUsername(mattermostUsername string) (*model.User, error)
	IsPluginAdmin(mattermostUserID string) (bool, error)
	UpdateStoredConfig(f func(*config.Config))
	Clean() error
}

// Dependencies contains all API dependencies
type Dependencies struct {
	Autofillers map[string]Autofiller
	PluginAPI
	Logger        bot.Logger
	Poster        bot.Poster
	RotationStore store.RotationStore
	ShiftStore    store.ShiftStore
	SkillsStore   store.SkillsStore
	UserStore     store.UserStore
}

type Config struct {
	*config.Config
	*Dependencies
}

type api struct {
	bot.Logger
	Config

	// set by `api.New`
	actingMattermostUserID string

	// use withActingUser or withActingUserExpanded to initialize.
	actingUser *User

	// use withKnownSkills to initialize.
	knownSkills store.IDMap

	// use withKnownRotations or withRotation(rotationID) to initialize, not expanded by default.
	knownRotations store.IDMap

	// use withMattermostUsers(usernames) or withUsers(mattermostUserIDs) to initialize, not expanded by default.
	users UserMap
}

func New(apiConfig Config, mattermostUserID string) API {
	return &api{
		Logger: apiConfig.Logger.With(bot.LogContext{
			"MattermostUserID": mattermostUserID,
		}),
		Config:                 apiConfig,
		actingMattermostUserID: mattermostUserID,
	}
}

func (api *api) MarkdownIndent(in, prefix string) string {
	lines := strings.Split(in, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
