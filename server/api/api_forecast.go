// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
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

type Forecaster interface {
	Guess(rotation *Rotation, startingShiftNumber, numShifts int, autofill bool) ([]*Shift, error)
	Forecast(rotation *Rotation, startingShiftNumber, numShifts, sampleSize int) (*Forecast, error)
}

func (api *api) Forecast(rotation *Rotation, startingShiftNumber, numShifts, sampleSize int) (*Forecast, error) {
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
		shifts, err = api.Guess(rotation, startingShiftNumber, numShifts, true)
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
				for _, need := range aerr.unfulfilledNeeds {
					f.NeedErrCounts[need.String()]++
				}

			case ErrFailedInsufficientForSize:
				f.CountErrFailedInsufficientForSize++
				for _, need := range aerr.unfulfilledNeeds {
					f.NeedErrCounts[need.String()]++
				}

			case ErrFailedSizeExceeded:
				f.CountErrFailedSizeExceeded++
			}

			f.ShiftErrCounts[aerr.shiftNumber]++
			continue GUESS
		}

		for n, shift := range shifts {
			for _, user := range shift.Users {
				if shift.MattermostUserIDs[user.MattermostUserID] != "" {
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
	}

	logger.Infof("Ran forecast for %s", MarkdownRotation(rotation))
	return f, nil
}

func (api *api) Guess(rotation *Rotation, startingShiftNumber int, numShifts int, autofill bool) ([]*Shift, error) {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	// logger := api.Logger.Timed().With(bot.LogContext{
	// 	"Location":       "api.Guess",
	// 	"ActingUsername": api.actingUser.MattermostUsername(),
	// 	"NumShifts":      numShifts,
	// 	"Autofill":       autofill,
	// 	"ShiftNumber":    startingShiftNumber,
	// 	"RotationID":     rotation.RotationID,
	// })

	prevUsers := rotation.Users
	rotation.Users = rotation.Users.Clone(true)
	defer func() {
		rotation.Users = prevUsers
	}()

	// Need to cache all users between iterations, to mock their last served dates.
	cachedUsers := UserMap{}
	for mattermostUserID, user := range rotation.Users {
		cachedUsers[mattermostUserID] = user
	}

	var shifts []*Shift
	for shiftNumber := startingShiftNumber; shiftNumber < startingShiftNumber+numShifts; shiftNumber++ {
		var shift *Shift
		shift, _, err := api.loadOrMakeOneShift(rotation, shiftNumber, autofill)
		if err != nil {
			if !autofill && err == store.ErrNotFound {
				shifts = append(shifts, nil)
				continue
			}
			return nil, err
		}

		// If the Shift was loaded, it might have brought some new users with
		// it. Replace them with cached Users as appropriate, they are more up-
		// to-date.
		for k := range shift.Users {
			if cachedUsers[k] != nil {
				shift.Users[k] = cachedUsers[k]
			}
		}

		if autofill && (shift.ShiftStatus == "" || shift.ShiftStatus == store.ShiftStatusOpen) {
			err = api.autofillShift(rotation, shiftNumber, shift, autofill)
			if err != nil {
				return nil, err
			}
		}

		// Update shift's users' last served counter, and update the cache in
		// case they were not there
		for k, u := range shift.Users {
			u.NextRotationShift[rotation.RotationID] = shiftNumber + 1 + rotation.Grace
			cachedUsers[k] = u
		}

		shifts = append(shifts, shift)
	}

	// logger.Debugf("Ran guess for %s", MarkdownRotation(rotation))
	return shifts, nil
}
