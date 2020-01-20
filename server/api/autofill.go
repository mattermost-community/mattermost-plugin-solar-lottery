// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

var ErrInsufficientForNeeds = errors.New("failed to satisfy needs, not enough skilled users available")
var ErrSizeExceeded = errors.New("failed to satisfy needs, exceeded rotation size")
var ErrInsufficientForSize = errors.New("failed to satisfy rotation size requirement")

type Autofiller interface {
	FillShift(rotation *Rotation, shiftNumber int, shift *Shift, logger bot.Logger) (UserMap, error)
}

type AutofillError struct {
	Err           error
	UnmetNeeds    store.Needs
	UnmetNeed     *store.Need
	UnmetCapacity int
	ShiftNumber   int
}

func (e AutofillError) Error() string {
	message := ""
	if e.UnmetCapacity > 0 {
		message = fmt.Sprintf("failed filling to capacity, missing %v", e.UnmetCapacity)
	}
	if e.UnmetNeed != nil {
		if message != "" {
			message += ", "
		}
		message += fmt.Sprintf("failed filling need %s", e.UnmetNeed.String())
	}
	if len(e.UnmetNeeds) > 0 {
		//TODO add message
	}
	if e.Err != nil {
		message = errors.WithMessage(e.Err, message).Error()
	}
	return message
}
