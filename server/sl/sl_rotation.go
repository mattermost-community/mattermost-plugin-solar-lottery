// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

var ErrMultipleResults = errors.New("multiple resolts found")

type RotationService interface {
	AddRotation(*Rotation) error
	ArchiveRotation(rotationID types.ID) (*Rotation, error)
	DebugDeleteRotation(rotationID types.ID) error
	LoadActiveRotations() (*types.IDIndex, error)
	LoadRotation(rotationID types.ID) (*Rotation, error)
	MakeRotation(rotationName string) (*Rotation, error)
	ResolveRotationName(string) (types.ID, error)
}

func (sl *sl) AddRotation(r *Rotation) error {
	var active *types.IDIndex
	err := sl.Setup(
		pushLogger("AddRotation", bot.LogContext{ctxRotationID: r.RotationID}),
		withActiveRotations(&active),
	)
	if err != nil {
		return err
	}
	defer sl.popLogger()

	if active.Contains(r.RotationID) {
		return ErrAlreadyExists
	}

	err = sl.Store.IDIndex(KeyActiveRotations).Set(r.RotationID)
	if err != nil {
		return err
	}
	active.Set(r.RotationID)
	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return err
	}

	sl.Infof("New rotation %s added", r.Markdown())
	return nil
}

func (sl *sl) LoadActiveRotations() (*types.IDIndex, error) {
	var active *types.IDIndex
	err := sl.Setup(withActiveRotations(&active))
	if err != nil {
		return nil, err
	}
	return active, nil
}

func (sl *sl) ResolveRotationName(pattern string) (types.ID, error) {
	var active *types.IDIndex
	err := sl.Setup(withActiveRotations(&active))
	if err != nil {
		return "", err
	}

	if active.Contains(types.ID(pattern)) {
		// Exact match
		return types.ID(pattern), nil
	}

	ids := []types.ID{}
	re, err := regexp.Compile(`.*` + pattern + `.*`)
	if err != nil {
		return "", err
	}

	for _, id := range active.IDs() {
		if re.MatchString(string(id)) {
			ids = append(ids, id)
		}
	}

	switch len(ids) {
	case 0:
		return "", kvstore.ErrNotFound
	case 1:
		return ids[0], nil
	}

	return "", errors.Errorf("ambiguous results: %v", ids)
}

func (sl *sl) ArchiveRotation(rotationID types.ID) (*Rotation, error) {
	var r *Rotation
	err := sl.Setup(
		pushLogger("ArchiveRotation", nil),
		withRotation(rotationID, &r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	r.IsArchived = true

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return nil, err
	}
	err = sl.Store.IDIndex(KeyActiveRotations).Delete(r.RotationID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to store rotation %s", r.RotationID)
	}

	sl.Infof("%s archived rotation %s.", sl.actingUser.Markdown(), r.Markdown())
	return r, nil
}

func (sl *sl) DebugDeleteRotation(rotationID types.ID) error {
	err := sl.Setup(pushLogger("DebugDeleteRotation", bot.LogContext{ctxRotationID: rotationID}))
	if err != nil {
		return err
	}
	defer sl.popLogger()

	err = sl.Store.Entity(KeyRotation).Delete(rotationID)
	if err != nil {
		return err
	}
	err = sl.Store.IDIndex(KeyActiveRotations).Delete(rotationID)
	if err != nil {
		return err
	}

	sl.Infof("%s deleted rotation %s.", sl.actingUser.Markdown(), rotationID)
	return nil
}

func (sl *sl) loadRotation(rotationID types.ID) (*Rotation, error) {
	var active *types.IDIndex
	err := sl.Setup(withActiveRotations(&active))
	if err != nil {
		return nil, err
	}

	if !active.Contains(rotationID) {
		return nil, errors.Errorf("rotationID %s not found", rotationID)
	}

	r := NewRotation()
	err = sl.Store.Entity(KeyRotation).Load(rotationID, r)
	if err != nil {
		return nil, err
	}
	r.init()

	return r, nil
}

func (sl *sl) MakeRotation(rotationName string) (*Rotation, error) {
	id, err := sl.Store.Entity(KeyRotation).NewID(rotationName)
	if err != nil {
		return nil, err
	}
	r := NewRotation()
	r.RotationID = id
	return r, nil
}

func (sl *sl) expandRotation(r *Rotation) error {
	if r.MattermostUserIDs.IsEmpty() ||
		(!r.users.IsEmpty() && r.users.Len() == r.MattermostUserIDs.Len()) {
		return nil
	}

	r.init()
	users, err := sl.loadStoredUsers(r.MattermostUserIDs)
	if err != nil {
		return err
	}
	err = sl.expandUsers(users)
	if err != nil {
		return err
	}
	r.users = users
	return nil
}

func (sl *sl) LoadRotation(rotationID types.ID) (*Rotation, error) {
	var r *Rotation
	err := sl.Setup(
		// pushLogger("LoadRotation", nil),
		withExpandedRotation(rotationID, &r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	return r, nil
}

func (sl *sl) updateRotation(r *Rotation, updatef func() error) error {
	err := sl.Setup(pushLogger("updateRotation", bot.LogContext{ctxRotationID: r.RotationID}))
	if err != nil {
		return err
	}
	defer sl.popLogger()

	err = updatef()
	if err != nil {
		return err
	}

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return err
	}

	sl.Debugf("%s updated rotation %s.", sl.actingUser.Markdown(), r.Markdown())
	return nil
}
