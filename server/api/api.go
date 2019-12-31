// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type API interface {
	bot.Logger
	Forecaster
	PluginAPI
	Rotations
	Shifts
	Skills
	Users
	UserActions
	Expander
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

type Shifts interface {
	CommitShift(*Rotation, int) error
	ListShifts(*Rotation, int, int) ([]*Shift, error)
	OpenShift(*Rotation, int) (*Shift, error)
	StartShift(*Rotation, int) error
	FinishShift(*Rotation, int) error
	DebugDeleteShift(*Rotation, int) error
	FillShift(*Rotation, int, bool) (shift *Shift, ready bool, whyNot string, before, added UserMap, err error)
}

type Users interface {
	GetActingUser() (*User, error)
	LoadMattermostUsers(mattermostUsernames string) (UserMap, error)
	LoadStoredUsers(mattermostUserIDs store.IDMap) (UserMap, error)
}

type UserActions interface {
	Qualify(mattermostUsernames, skillName string, level Level) error
	JoinRotation(mattermostUsernames string, rotation *Rotation, graceShifts int) (added UserMap, err error)
	LeaveRotation(mattermostUsernames string, rotation *Rotation) (deleted UserMap, err error)
	Disqualify(mattermostUsernames, skillName string) error
	JoinShift(mattermostUsernames string, rotation *Rotation, shiftNumber int) error
	AddEvent(mattermostUsernames string, event store.Event) error
	DeleteEvents(mattermostUsernames string, startDate, endDate string) error
}

type Forecaster interface {
	Guess(rotation *Rotation, startingShiftNumber, numShifts int, autofill bool) ([]*Shift, error)
	ForecastRotation(rotation *Rotation, startingShiftNumber, numShifts, sampleSize int) (*Forecast, error)
	ForecastUser(mattermostUsername string, rotation *Rotation, numShifts, sampleSize int) ([]float64, error)
}

type Skills interface {
	ListSkills() (store.IDMap, error)
	AddSkill(string) error
	DeleteSkill(string) error
}

type Expander interface {
	ExpandUserMap(UserMap) error
	ExpandUser(*User) error
	ExpandShift(*Shift) error
	ExpandRotation(*Rotation) error
}

type PluginAPI interface {
	GetMattermostUser(mattermostUserID string) (*model.User, error)
	GetMattermostUserByUsername(mattermostUsername string) (*model.User, error)
	IsPluginAdmin(mattermostUserID string) (bool, error)
}

// Dependencies contains all API dependencies
type Dependencies struct {
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

	// use withRotation(rotationID) to initialize, not expanded by default.
	// rotation *Rotation

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

type filterf func(*api) error

func (api *api) Filter(filters ...filterf) error {
	for _, filter := range filters {
		err := filter(api)
		if err != nil {
			return err
		}
	}
	return nil
}
