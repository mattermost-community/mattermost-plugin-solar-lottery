// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/pkg/errors"
)

const DayDuration = time.Hour * 24
const WeekDuration = DayDuration * 7
const DateFormat = "2006-01-02"

type Shifts interface {
	ListShifts(*Rotation, int, int) ([]*Shift, error)
	OpenShift(*Rotation, int) (*Shift, error)
	StartShift(*Rotation, int) (*Shift, error)
	FinishShift(*Rotation, int) (*Shift, error)
	DebugDeleteShift(*Rotation, int) error
	FillShift(*Rotation, int) (*Shift, UserMap, error)
	IsShiftReady(rotation *Rotation, shiftNumber int) (shift *Shift, ready bool, whyNot string, err error)
}

type Shift struct {
	*store.Shift
	StartTime time.Time
	EndTime   time.Time

	ShiftNumber  int
	RotationName string
}

func (shift *Shift) Clone(deep bool) *Shift {
	newShift := *shift
	if deep {
		newShift.Shift = &(*shift.Shift)
	}
	return &newShift
}

func (shift Shift) MarkdownBullets(rotation *Rotation) string {
	out := fmt.Sprintf("- %s\n", shift.Markdown())
	out += fmt.Sprintf("  - Status: **%s**\n", shift.Status)
	out += fmt.Sprintf("  - Users: **%v**\n", len(shift.MattermostUserIDs))
	for _, user := range rotation.ShiftUsers(&shift) {
		out += fmt.Sprintf("    - %s\n", user.MarkdownWithSkills())
	}
	return out
}

func (shift Shift) Markdown() string {
	return fmt.Sprintf("%s#%v", shift.RotationName, shift.ShiftNumber)
}

func (rotation *Rotation) makeShift(shiftNumber int) (*Shift, error) {
	start, end, err := rotation.ShiftDatesForNumber(shiftNumber)
	if err != nil {
		return nil, err
	}
	return &Shift{
		Shift:        store.NewShift(start.Format(DateFormat), end.Format(DateFormat), nil),
		StartTime:    start,
		EndTime:      end,
		RotationName: rotation.Name,
		ShiftNumber:  shiftNumber,
	}, nil
}

func (sl *solarLottery) loadShift(rotation *Rotation, shiftNumber int) (*Shift, error) {
	shift, err := rotation.makeShift(shiftNumber)
	if err != nil {
		return nil, err
	}
	s, err := sl.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	if err != nil {
		return nil, err
	}
	shift.Shift = s

	shift.StartTime, shift.EndTime, err = ParseDatePair(s.Start, s.End)
	if err != nil {
		return nil, err
	}

	shift.RotationName = rotation.Name
	shift.ShiftNumber = shiftNumber
	return shift, nil
}

// Returns an un-expanded shift - will be populated with Users from rotation
func (sl *solarLottery) getShiftForGuess(rotation *Rotation, shiftNumber int) (*Shift, bool, error) {
	start, end, err := rotation.ShiftDatesForNumber(shiftNumber)
	if err != nil {
		return nil, false, err
	}

	var shift *Shift
	created := false
	storedShift, err := sl.ShiftStore.LoadShift(rotation.RotationID, shiftNumber)
	switch err {
	case nil:
		shift = &Shift{
			Shift: storedShift,
		}

	case store.ErrNotFound:
		shift, err = rotation.makeShift(shiftNumber)
		if err != nil {
			return nil, false, err
		}
		created = true

	default:
		return nil, false, err
	}

	if shift.Start != start.Format(DateFormat) || shift.End != end.Format(DateFormat) {
		return nil, false, errors.Errorf("loaded shift has wrong dates %v-%v, expected %v-%v",
			shift.Start, shift.End, start, end)
	}

	shift.RotationName = rotation.Name
	shift.ShiftNumber = shiftNumber
	return shift, created, nil
}
