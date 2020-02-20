// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package autofill

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
)

var ErrInsufficientForNeeds = errors.New("failed to satisfy needs, not enough skilled users available")
var ErrSizeExceeded = errors.New("failed to satisfy needs, exceeded rotation size")
var ErrInsufficientForSize = errors.New("failed to satisfy rotation size requirement")

type Error struct {
	Err           error
	UnmetNeeds    []*sl.Need
	UnmetNeed     *sl.Need
	UnmetCapacity int
	ShiftNumber   int
}

func (e Error) Error() string {
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
