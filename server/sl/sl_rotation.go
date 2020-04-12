// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (sl *sl) AddRotation(r *Rotation) error {
	active := types.NewIDSet()
	err := sl.Setup(
		pushAPILogger("AddRotation", r),
		withLoadIDIndex(KeyActiveRotations, active),
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

func (sl *sl) LoadActiveRotations() (*types.IDSet, error) {
	active := types.NewIDSet()
	err := sl.Setup(withLoadIDIndex(KeyActiveRotations, active))
	if err != nil {
		return nil, err
	}
	return active, nil
}

func (sl *sl) ResolveRotationName(pattern string) (types.ID, error) {
	active := types.NewIDSet()
	err := sl.Setup(withLoadIDIndex(KeyActiveRotations, active))
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
	r := NewRotation()
	err := sl.Setup(
		pushAPILogger("ArchiveRotation", rotationID),
		withLoadRotation(&rotationID, r),
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
	err := sl.Setup(pushAPILogger("DebugDeleteRotation", rotationID))
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

func (sl *sl) LoadRotation(rotationID types.ID) (*Rotation, error) {
	r := NewRotation()
	err := sl.Setup(withExpandedRotation(&rotationID, r))
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	return r, nil
}

func (sl *sl) UpdateRotation(rotationID types.ID, updatef func(*Rotation) error) (*Rotation, error) {
	r := NewRotation()
	err := sl.Setup(withLoadRotation(&rotationID, r))
	if err != nil {
		return nil, err
	}

	err = updatef(r)
	if err != nil {
		return nil, err
	}

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return nil, err
	}

	sl.Debugf("%s updated rotation %s.", sl.actingUser.Markdown(), r.Markdown())
	return r, nil
}

func (sl *sl) MakeRotation(rotationName string) (*Rotation, error) {
	id, err := sl.Store.Entity(KeyRotation).NewID(rotationName)
	if err != nil {
		return nil, err
	}
	r := NewRotation()
	r.RotationID = id
	r.loaded = true
	return r, nil
}

func (sl *sl) loadRotation(rotationID types.ID) (*Rotation, error) {
	active := types.NewIDSet()
	err := sl.Setup(withLoadIDIndex(KeyActiveRotations, active))
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
	r.loaded = true

	return r, nil
}

func (sl *sl) expandRotationUsers(r *Rotation) error {
	users, err := sl.LoadUsers(r.MattermostUserIDs)
	if err != nil {
		return err
	}
	r.Users = users
	return nil
}

func (sl *sl) expandRotationTasks(r *Rotation) error {
	if r.Tasks != nil { // && r.inProgress != nil {
		return nil
	}

	tasks, err := sl.LoadTasks(r.TaskIDs)
	if err != nil {
		return err
	}
	r.Tasks = tasks
	return nil
}

func (sl *sl) taskFiller(r *Rotation) (TaskFiller, error) {
	f, ok := sl.TaskFillers[r.FillerType]
	if !ok {
		return nil, errors.Wrapf(kvstore.ErrNotFound, "filler type %s", r.FillerType)
	}
	return f, nil
}
