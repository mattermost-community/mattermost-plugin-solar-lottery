// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type TaskFiller interface {
	FillTask(rotation *Rotation, task *Task, forTime types.Time, logger bot.Logger) (*Users, error)
}

var ErrFillInsufficient = errors.New("insufficient")
var ErrFillLimit = errors.New("limit violated")

type FillError struct {
	Err        error
	FailedNeed *Need
	UnmetNeeds *Needs
	TaskID     types.ID
}

func (e FillError) Error() string {
	message := fmt.Sprintf("failed to fill %s", e.TaskID)
	if e.FailedNeed != nil {
		if message != "" {
			message += ", "
		}
		message += fmt.Sprintf("filling need %s", e.FailedNeed.Markdown())
	}
	if !e.UnmetNeeds.IsEmpty() {
		if message != "" {
			message += ", "
		}
		message += fmt.Sprintf("unfilled needs %s", e.UnmetNeeds.Markdown())
	}
	if e.Err != nil {
		message = errors.WithMessage(e.Err, message).Error()
	}
	return message
}
