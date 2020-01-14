// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Forecast struct {
	StartingShift                      int
	NumShifts                          int
	SampleSize                         int
	CountErrFailedInsufficientForNeeds int
	CountErrFailedInsufficientForSize  int
	CountErrFailedSizeExceeded         int
	NeedErrCounts                      map[string]int
	ShiftErrCounts                     []int

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
		var aerr autofillError
		if err != nil {
			var ok bool
			aerr, ok = err.(autofillError)
			if !ok {
				return nil, err
			}

			switch aerr.orig {
			case ErrFailedInsufficientForNeeds:
				f.CountErrFailedInsufficientForNeeds++
				for _, need := range aerr.causeUnmetNeeds {
					f.NeedErrCounts[need.String()]++
				}

			case ErrFailedInsufficientForSize:
				f.CountErrFailedInsufficientForSize++
				for _, need := range aerr.causeUnmetNeeds {
					f.NeedErrCounts[need.String()]++
				}

			case ErrFailedSizeExceeded:
				f.CountErrFailedSizeExceeded++
			}

			f.ShiftErrCounts[aerr.shiftNumber]++
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

	logger.Infof("Ran forecast for %s", api.MarkdownRotation(rotation))
	return f, nil
}
