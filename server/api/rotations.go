// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

var ErrMultipleResults = errors.New("multiple resolts found")

func (api *api) AddRotation(rotation *Rotation) error {
	err := api.Filter(
		withActingUserExpanded,
		withKnownRotations,
		rotation.init,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.AddRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})
	_, ok := api.knownRotations[rotation.RotationID]
	if ok {
		return ErrAlreadyExists
	}

	api.knownRotations[rotation.RotationID] = rotation.Name
	err = api.RotationStore.StoreKnownRotations(api.knownRotations)
	if err != nil {
		return err
	}
	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}
	logger.Infof("New rotation %s added", rotation.Markdown())
	return nil
}

func (api *api) LoadKnownRotations() (store.IDMap, error) {
	err := api.Filter(
		withActingUser,
		withKnownRotations,
	)
	if err != nil {
		return nil, err
	}
	return api.knownRotations, nil
}

func (api *api) ResolveRotationName(namePattern string) ([]string, error) {
	err := api.Filter(
		withKnownRotations,
	)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	re, err := regexp.Compile(`.*` + namePattern + `.*`)
	if err != nil {
		return nil, err
	}
	for id, name := range api.knownRotations {
		if re.MatchString(name) {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return nil, store.ErrNotFound
	}
	return ids, nil
}

func (api *api) ArchiveRotation(rotation *Rotation) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.ArchiveRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})

	rotation.Rotation.IsArchived = true

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}
	delete(api.knownRotations, rotation.RotationID)
	err = api.RotationStore.StoreKnownRotations(api.knownRotations)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", rotation.RotationID)
	}

	logger.Infof("%s archived rotation %s.", api.actingUser.Markdown(), rotation.Markdown())
	return nil
}

func (api *api) DebugDeleteRotation(rotationID string) error {
	err := api.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.DebugDeleteRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotationID,
	})

	err = api.RotationStore.DeleteRotation(rotationID)
	if err != nil {
		return err
	}
	delete(api.knownRotations, rotationID)
	err = api.RotationStore.StoreKnownRotations(api.knownRotations)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", rotationID)
	}

	logger.Infof("%s deleted rotation %s.", api.actingUser.Markdown(), rotationID)
	return nil
}

func (api *api) LoadRotation(rotationID string) (*Rotation, error) {
	err := api.Filter(
		withKnownRotations,
	)
	if err != nil {
		return nil, err
	}

	_, ok := api.knownRotations[rotationID]
	if !ok {
		return nil, errors.Errorf("rotationID %s not found", rotationID)
	}

	storedRotation, err := api.RotationStore.LoadRotation(rotationID)
	rotation := &Rotation{
		Rotation: storedRotation,
	}
	err = rotation.init(api)
	if err != nil {
		return nil, err
	}

	return rotation, nil
}

func (api *api) MakeRotation(rotationName string) (*Rotation, error) {
	id := ""
	for i := 0; i < 5; i++ {
		tryId := rotationName + "-" + model.NewId()[:7]
		if len(api.knownRotations) == 0 || api.knownRotations[tryId] == "" {
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

func (api *api) UpdateRotation(rotation *Rotation, updatef func(*Rotation) error) error {
	err := api.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":       "api.UpdateRotation",
		"ActingUsername": api.actingUser.MattermostUsername(),
		"RotationID":     rotation.RotationID,
	})

	err = updatef(rotation)
	if err != nil {
		return err
	}

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return err
	}

	logger.Infof("%s updated rotation %s.", api.actingUser.Markdown(), rotation.Markdown())
	return nil
}
