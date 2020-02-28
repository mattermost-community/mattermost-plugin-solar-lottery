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

type Rotations interface {
	ExpandRotation(*Rotation) error
	JoinRotation(*Rotation, UserMap, types.Time) (added UserMap, err error)
	LeaveRotation(*Rotation, UserMap) (deleted UserMap, err error)
	AddRotation(*Rotation) error
	ArchiveRotation(*Rotation) error
	DebugDeleteRotation(types.ID) error
	LoadRotation(types.ID) (*Rotation, error)
	MakeRotation(rotationName string) (*Rotation, error)
	ResolveRotation(string) (types.ID, error)
	UpdateRotation(*Rotation, func(*Rotation) error) error
	LoadActiveRotations() (*types.IDIndex, error)
}

func (sl *sl) AddRotation(r *Rotation) error {
	err := sl.Filter(
		withActingUserExpanded,
		withActiveRotations,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":   "AddRotation",
		"ActingUser": sl.actingUser.MattermostUsername(),
		"RotationID": r.RotationID,
	})
	if sl.activeRotations.Contains(r.RotationID) {
		return ErrAlreadyExists
	}

	err = sl.Store.IDIndex(KeyActiveRotations).Set(r.RotationID)
	if err != nil {
		return err
	}
	sl.activeRotations.Set(r.RotationID)
	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return err
	}

	logger.Infof("New rotation %s added", r.Markdown())
	return nil
}

func (sl *sl) LoadActiveRotations() (*types.IDIndex, error) {
	err := sl.Filter(
		withActingUser,
		withActiveRotations,
	)
	if err != nil {
		return nil, err
	}
	return sl.activeRotations, nil
}

func (sl *sl) ResolveRotation(pattern string) (types.ID, error) {
	err := sl.Filter(
		withActiveRotations,
	)
	if err != nil {
		return "", err
	}

	if sl.activeRotations.Contains(types.ID(pattern)) {
		// Exact match
		return types.ID(pattern), nil
	}

	ids := []types.ID{}
	re, err := regexp.Compile(`.*` + pattern + `.*`)
	if err != nil {
		return "", err
	}

	for _, id := range sl.activeRotations.IDs() {
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

func (sl *sl) ArchiveRotation(r *Rotation) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":   "sl.ArchiveRotation",
		"ActingUser": sl.actingUser.MattermostUsername(),
		"RotationID": r.RotationID,
	})

	r.IsArchived = true

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return err
	}
	err = sl.Store.IDIndex(KeyActiveRotations).Delete(r.RotationID)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", r.RotationID)
	}
	sl.activeRotations.Delete(r.RotationID)

	logger.Infof("%s archived rotation %s.", sl.actingUser.Markdown(), r.Markdown())
	return nil
}

func (sl *sl) DebugDeleteRotation(rotationID types.ID) error {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":   "DebugDeleteRotation",
		"ActingUser": sl.actingUser.MattermostUsername(),
		"RotationID": rotationID,
	})

	err = sl.Store.Entity(KeyRotation).Delete(rotationID)
	if err != nil {
		return err
	}
	err = sl.Store.IDIndex(KeyActiveRotations).Delete(rotationID)
	if err != nil {
		return err
	}
	sl.activeRotations.Delete(rotationID)

	logger.Infof("%s deleted rotation %s.", sl.actingUser.Markdown(), rotationID)
	return nil
}

func (sl *sl) LoadRotation(rotationID types.ID) (*Rotation, error) {
	err := sl.Filter(
		withActiveRotations,
	)
	if err != nil {
		return nil, err
	}

	if !sl.activeRotations.Contains(rotationID) {
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

func (sl *sl) ExpandRotation(r *Rotation) error {
	if r.MattermostUserIDs != nil && len(r.users) != r.MattermostUserIDs.Len() {
		users, err := sl.loadStoredUsers(r.MattermostUserIDs)
		if err != nil {
			return err
		}
		err = sl.ExpandUserMap(users)
		if err != nil {
			return err
		}
		r.users = users
	}

	return nil
}

func (sl *sl) UpdateRotation(r *Rotation, updatef func(*Rotation) error) error {
	err := sl.Filter(
		withActingUserExpanded,
		withRotationExpanded(r),
	)
	if err != nil {
		return err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":   "UpdateRotation",
		"ActingUser": sl.actingUser.MattermostUsername(),
		"RotationID": r.RotationID,
	})

	err = updatef(r)
	if err != nil {
		return err
	}

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return err
	}

	logger.Infof("%s updated rotation %s.", sl.actingUser.Markdown(), r.Markdown())
	return nil
}

func (sl *sl) JoinRotation(r *Rotation, users UserMap, starting types.Time) (UserMap, error) {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":   "Join",
		"ActingUser": sl.actingUser.MattermostUsername(),
		"RotationID": r.RotationID,
		"Users":      users.String(),
		"Starting":   starting,
	})

	added := UserMap{}
	for _, user := range users {
		if r.MattermostUserIDs.Contains(user.MattermostUserID) {
			logger.Debugf("%s is already in rotation %s.",
				added.MarkdownWithSkills(), r.Markdown())
			continue
		}

		// A new person may be given some slack - setting starting in the
		// future all but guarantees they won't be selected until then.
		user.LastServed.Set(r.RotationID, starting.Unix())

		user, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return added, err
		}

		if r.MattermostUserIDs == nil {
			r.MattermostUserIDs = types.NewIDIndex()
		}
		r.MattermostUserIDs.Set(user.MattermostUserID)
		sl.messageWelcomeToRotation(user, r)
		added[user.MattermostUserID] = user
	}

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return added, errors.WithMessagef(err, "failed to store rotation %s", r.RotationID)
	}
	logger.Infof("%s added %s to %s.",
		sl.actingUser.Markdown(), added.MarkdownWithSkills(), r.Markdown())
	return added, nil
}

func (sl *sl) LeaveRotation(r *Rotation, users UserMap) (UserMap, error) {
	err := sl.Filter(
		withActingUserExpanded,
	)
	if err != nil {
		return nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "LeaveFromRotation",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"Users":          users.String(),
		"RotationID":     r.RotationID,
	})

	deleted := UserMap{}
	for _, user := range users {
		if !r.MattermostUserIDs.Contains(user.MattermostUserID) {
			logger.Debugf("%s is not found in rotation %s", user.Markdown(), r.Markdown())
			continue
		}

		user.LastServed.Delete(r.RotationID)
		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return deleted, err
		}
		r.MattermostUserIDs.Delete(user.MattermostUserID)
		if len(r.users) > 0 {
			delete(r.users, user.MattermostUserID)
		}
		sl.messageLeftRotation(user, r)
		deleted[user.MattermostUserID] = user
	}

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return deleted, err
	}

	logger.Infof("%s removed from %s.", deleted.Markdown(), r.Markdown())
	return deleted, nil
}
