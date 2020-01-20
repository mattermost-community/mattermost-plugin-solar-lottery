// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

type Forecast struct {
	StartingShift                int
	NumShifts                    int
	SampleSize                   int
	CountErrInsufficientForNeeds int
	CountErrInsufficientForSize  int
	CountErrSizeExceeded         int
	NeedErrCounts                map[string]int
	ShiftErrCounts               []int

	// for each user (MattermostUserID) contains NumShifts probabilities, of
	// the user being selected into the respective shift number. This is based
	// on successful guesses only.
	UserShiftCounts map[string][]int
	UserCounts      map[string]int
}

func (api *api) ForecastRotation(rotation *Rotation, startingShiftNumber, numShifts, sampleSize int) (*Forecast, error) {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.Forecast",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"NumShifts":           numShifts,
		"StartingShiftNumber": startingShiftNumber,
		"RotationID":          rotation.RotationID,
	})

	f := &Forecast{
		StartingShift:   startingShiftNumber,
		NumShifts:       numShifts,
		SampleSize:      sampleSize,
		NeedErrCounts:   map[string]int{},
		ShiftErrCounts:  make([]int, numShifts),
		UserShiftCounts: map[string][]int{},
		UserCounts:      map[string]int{},
	}

GUESS:
	for i := 0; i < sampleSize; i++ {
		var shifts []*Shift
		// Guess' logs are too verbose - suppress
		prevLogger := api.Logger
		api.Logger = &bot.NilLogger{}
		shifts, err = api.Guess(rotation, startingShiftNumber, numShifts)
		api.Logger = prevLogger
		var aerr *AutofillError
		if err != nil {
			var ok bool
			aerr, ok = err.(*AutofillError)
			if !ok {
				return nil, err
			}

			for _, need := range aerr.UnmetNeeds {
				f.NeedErrCounts[need.String()]++
			}

			switch aerr.Err {
			case ErrInsufficientForNeeds:
				f.CountErrInsufficientForNeeds++

			case ErrInsufficientForSize:
				f.CountErrInsufficientForSize++

			case ErrSizeExceeded:
				f.CountErrSizeExceeded++
			}

			f.ShiftErrCounts[aerr.ShiftNumber]++
			continue GUESS
		}

		for n, shift := range shifts {
			for _, user := range rotation.ShiftUsers(shift) {
				sc := f.UserShiftCounts[user.MattermostUsername()]
				if sc == nil {
					sc = make([]int, numShifts)
				}
				sc[n]++
				f.UserShiftCounts[user.MattermostUsername()] = sc
				f.UserCounts[user.MattermostUsername()]++
			}
		}
	}

	logger.Infof("Ran forecast for %s", rotation.Markdown())
	return f, nil
}

func (api *api) ForecastUser(mattermostUsername string, rotation *Rotation, numShifts, sampleSize int, now time.Time) ([]float64, error) {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsername),
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	if len(api.users) != 1 {
		return nil, errors.New("must provide a single user name")
	}

	var user *User
	for _, u := range api.users {
		user = u
		break
	}

	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.Forecast",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"NumShifts":      numShifts,
		"Username":       mattermostUsername,
		"RotationID":     rotation.RotationID,
	})

	shiftNumber, err := rotation.ShiftNumberForTime(now)
	if err != nil {
		return nil, err
	}
	shiftNumber++ // start with the next shift, or 0 if -1 was returned

	shiftCounts := make([]float64, numShifts)

GUESS:
	for i := 0; i < sampleSize; i++ {
		var shifts []*Shift
		prevLogger := api.Logger
		api.Logger = &bot.NilLogger{}
		shifts, err = api.Guess(rotation, shiftNumber, numShifts)
		api.Logger = prevLogger
		if err != nil {
			continue GUESS
		}

		for n, shift := range shifts {
			if shift.MattermostUserIDs[user.MattermostUserID] != "" {
				shiftCounts[n]++
			}
		}
	}

	expectedServed := []float64{}
	var cumulative float64
	for _, c := range shiftCounts {
		cumulative += c
		expectedServed = append(expectedServed, cumulative/float64(sampleSize))
	}

	logger.Infof("Ran forecast for %s, user %s", rotation.Markdown(), user.Markdown())
	return expectedServed, nil
}
