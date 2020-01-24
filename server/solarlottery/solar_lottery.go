// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type SolarLottery interface {
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

type Autofiller interface {
	FillShift(rotation *Rotation, shiftNumber int, shift *Shift, logger bot.Logger) (UserMap, error)
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

type solarLottery struct {
	bot.Logger
	Config

	// set by `sl.New`
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

var _ SolarLottery = (*solarLottery)(nil)

func New(apiConfig Config, mattermostUserID string) SolarLottery {
	return &solarLottery{
		Logger: apiConfig.Logger.With(bot.LogContext{
			"MattermostUserID": mattermostUserID,
		}),
		Config:                 apiConfig,
		actingMattermostUserID: mattermostUserID,
	}
}
