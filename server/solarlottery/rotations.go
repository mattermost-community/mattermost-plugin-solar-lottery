// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrMultipleResults = errors.New("multiple resolts found")

type Rotations interface {
	AddRotation(*Rotation) error
	ArchiveRotation(*Rotation) error
	DebugDeleteRotation(string) error
	LoadActiveRotations() (store.IDMap, error)
	LoadRotation(string) (*Rotation, error)
	MakeRotation(rotationName string) (*Rotation, error)
	ResolveRotationName(namePattern string) ([]string, error)
	UpdateRotation(*Rotation, func(*Rotation) error) error
}

func (sl *solarLottery) AddRotation(rotation *Rotation) error {
	err := sl.Filter(
		withActingUserExpanded,
		withActiveRotations,
		rotation.init,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.AddRotation",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})
	_, ok := sl.activeRotations[rotation.RotationID]
	if ok {
		return ErrAlreadyExists
	}

	sl.activeRotations[rotation.RotationID] = rotation.Name
	err = sl.RotationStore.StoreActiveRotations(sl.activeRotations)
	if err != nil {
		return err
	}
	err = sl.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}
	logger.Infof("New rotation %s added", rotation.Markdown())
	return nil
}

func (sl *solarLottery) LoadActiveRotations() (store.IDMap, error) {
	err := sl.Filter(
		withActingUser,
		withActiveRotations,
	)
	if err != nil {
		return nil, err
	}
	return sl.activeRotations, nil
}

func (sl *solarLottery) ResolveRotationName(namePattern string) ([]string, error) {
	err := sl.Filter(
		withActiveRotations,
	)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	re, err := regexp.Compile(`.*` + namePattern + `.*`)
	if err != nil {
		return nil, err
	}
	for id, name := range sl.activeRotations {
		if re.MatchString(name) {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return nil, store.ErrNotFound
	}
	return ids, nil
}

func (sl *solarLottery) ArchiveRotation(rotation *Rotation) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.ArchiveRotation",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})

	rotation.Rotation.IsArchived = true

	err = sl.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}
	delete(sl.activeRotations, rotation.RotationID)
	err = sl.RotationStore.StoreActiveRotations(sl.activeRotations)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", rotation.RotationID)
	}

	logger.Infof("%s archived rotation %s.", sl.actingUser.Markdown(), rotation.Markdown())
	return nil
}

func (sl *solarLottery) DebugDeleteRotation(rotationID string) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.DebugDeleteRotation",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotationID,
	})

	err = sl.RotationStore.DeleteRotation(rotationID)
	if err != nil {
		return err
	}
	delete(sl.activeRotations, rotationID)
	err = sl.RotationStore.StoreActiveRotations(sl.activeRotations)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", rotationID)
	}

	logger.Infof("%s deleted rotation %s.", sl.actingUser.Markdown(), rotationID)
	return nil
}

func (sl *solarLottery) LoadRotation(rotationID string) (*Rotation, error) {
	err := sl.Filter(
		withActiveRotations,
	)
	if err != nil {
		return nil, err
	}

	_, ok := sl.activeRotations[rotationID]
	if !ok {
		return nil, errors.Errorf("rotationID %s not found", rotationID)
	}

	storedRotation, err := sl.RotationStore.LoadRotation(rotationID)
	rotation := &Rotation{
		Rotation: storedRotation,
	}
	err = rotation.init(sl)
	if err != nil {
		return nil, err
	}

	return rotation, nil
}

func (sl *solarLottery) MakeRotation(rotationName string) (*Rotation, error) {
	id := ""
	for i := 0; i < 5; i++ {
		tryId := rotationName + "-" + model.NewId()[:7]
		if len(sl.activeRotations) == 0 || sl.activeRotations[tryId] == "" {
			id = tryId
			break
		}
	}
	if id == "" {
		return nil, errors.New("Failed to generate unique rotation ID")
	}

	rotation := &Rotation{
		Rotation: store.NewRotation(rotationName),
	}
	rotation.RotationID = id
	return rotation, nil
}

func (sl *solarLottery) UpdateRotation(rotation *Rotation, updatef func(*Rotation) error) error {
	err := sl.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.UpdateRotation",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})

	err = updatef(rotation)
	if err != nil {
		return err
	}

	err = sl.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}

	logger.Infof("%s updated rotation %s.", sl.actingUser.Markdown(), rotation.Markdown())
	return nil
}
