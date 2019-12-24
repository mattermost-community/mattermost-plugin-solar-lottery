// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Rotations interface {
	MakeRotation(rotationName string) (*Rotation, error)
	AddRotation(*Rotation) error
	AddRotationUsers(rotation *Rotation, mattermostUsernames string, graceShifts int) (added UserMap, err error)
	ArchiveRotation(*Rotation) error
	DeleteRotationUsers(rotation *Rotation, mattermostUsernames string) (deleted UserMap, err error)
	LoadKnownRotations() (store.IDMap, error)
	LoadRotation(string) (*Rotation, error)
	DebugDeleteRotation(string) error
	ResolveRotationName(namePattern string) ([]string, error)
	UpdateRotation(*Rotation, func(*Rotation) error) error
}

var ErrRotationAlreadyExists = errors.New("rotation already exists")

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

	return &Rotation{
		Rotation: &store.Rotation{
			RotationID: id,
			Name:       rotationName,
		},
		Users: UserMap{},
	}, nil
}

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
		return ErrRotationAlreadyExists
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
	logger.Infof("New rotation %s added.", MarkdownRotationWithDetails(rotation))
	return nil
}

func (api *api) AddRotationUsers(rotation *Rotation, mattermostUsernames string, graceShifts int) (UserMap, error) {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.AddRotationUsers",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"RotationID":          rotation.RotationID,
		"MattermostUsernames": mattermostUsernames,
		"GraceShifts":         graceShifts,
	})

	shiftNumber, _ := rotation.ShiftNumberForTime(time.Now())
	added := UserMap{}
	for _, user := range api.users {
		if len(rotation.MattermostUserIDs[user.MattermostUserID]) != 0 {
			logger.Debugf("%s is already in rotation %s.",
				MarkdownUserMapWithSkills(added), MarkdownRotation(rotation))
			continue
		}

		// A new person may be given some slack - setting LastShiftNumber in the
		// future guarantees they won't be selected until then.
		user.Rotations[rotation.RotationID] = shiftNumber + graceShifts

		user, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return added, err
		}

		if rotation.MattermostUserIDs == nil {
			rotation.MattermostUserIDs = store.IDMap{}
		}

		rotation.MattermostUserIDs[user.MattermostUserID] = user.MattermostUserID
		api.messageWelcomeToRotation(user, rotation)
		added[user.MattermostUserID] = user
	}

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return added, errors.WithMessagef(err, "failed to store rotation %s", rotation.RotationID)
	}
	logger.Infof("%s added %s to %s.",
		MarkdownUser(api.actingUser), MarkdownUserMapWithSkills(added), MarkdownRotation(rotation))
	return added, nil
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

	logger.Infof("%s archived rotation %s.", MarkdownUser(api.actingUser), MarkdownRotation(rotation))
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

	logger.Infof("%s deleted rotation %s.", MarkdownUser(api.actingUser), rotationID)
	return nil
}

func (api *api) DeleteRotationUsers(rotation *Rotation, mattermostUsernames string) (UserMap, error) {
	err := api.Filter(
		withActingUserExpanded,
		withMattermostUsersExpanded(mattermostUsernames),
	)
	if err != nil {
		return nil, err
	}
	logger := api.Logger.Timed().With(bot.LogContext{
		"Location":            "api.DeleteUsersFromRotation",
		"ActingUsername":      api.actingUser.MattermostUsername(),
		"MattermostUsernames": mattermostUsernames,
		"RotationID":          rotation.RotationID,
	})

	deleted := UserMap{}
	for _, user := range api.users {
		_, ok := rotation.MattermostUserIDs[user.MattermostUserID]
		if !ok {
			logger.Debugf("%s is not found in rotation %s", MarkdownUser(user), MarkdownRotation(rotation))
			continue
		}

		delete(user.Rotations, rotation.RotationID)
		_, err = api.storeUserWelcomeNew(user)
		if err != nil {
			return deleted, err
		}
		delete(rotation.MattermostUserIDs, user.MattermostUserID)
		if len(rotation.Users) > 0 {
			delete(rotation.Users, user.MattermostUserID)
		}
		api.messageLeftRotation(user, rotation)
		deleted[user.MattermostUserID] = user
	}

	err = api.RotationStore.StoreRotation(rotation.Rotation)
	if err != nil {
		return deleted, err
	}

	logger.Infof("%s removed from %s.", MarkdownUserMap(deleted), MarkdownRotation(rotation))
	return deleted, nil
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
	return ids, nil
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

	logger.Infof("%s updated rotation %s.", MarkdownUser(api.actingUser), MarkdownRotationWithDetails(rotation))
	return nil
}
