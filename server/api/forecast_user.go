// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

func (api *api) ForecastUser(mattermostUsername string, rotation *Rotation, numShifts, sampleSize int) ([]float64, error) {
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

	shiftNumber, err := rotation.ShiftNumberForTime(time.Now())
	if err != nil {
		return nil, err
	}
	shiftNumber++ // start with the next shift

	shiftCounts := make([]float64, numShifts)

GUESS:
	for i := 0; i < sampleSize; i++ {
		var shifts []*Shift
		prevLogger := api.Logger
		api.Logger = &bot.NilLogger{}
		shifts, err = api.Guess(rotation, shiftNumber, numShifts, true)
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

	logger.Infof("Ran forecast for %s, user %s", MarkdownRotation(rotation), MarkdownUser(user))
	return expectedServed, nil
}
