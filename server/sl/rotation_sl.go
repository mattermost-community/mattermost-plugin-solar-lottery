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
	DebugDeleteRotation(string) error
	LoadRotation(string) (*Rotation, error)
	MakeRotation(rotationName string) (*Rotation, error)
	ResolveRotation(string) (string, error)
	UpdateRotation(*Rotation, func(*Rotation) error) error
	LoadActiveRotations() (*types.Set, error)
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

	sl.activeRotations.Add(r.RotationID)
	err = sl.Store.Index(KeyActiveRotations).Store(sl.activeRotations)
	if err != nil {
		return err
	}
	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return err
	}
	logger.Infof("New rotation %s added", r.Markdown())
	return nil
}

func (sl *sl) LoadActiveRotations() (*types.Set, error) {
	err := sl.Filter(
		withActingUser,
		withActiveRotations,
	)
	if err != nil {
		return nil, err
	}
	return sl.activeRotations, nil
}

func (sl *sl) ResolveRotation(pattern string) (string, error) {
	err := sl.Filter(
		withActiveRotations,
	)
	if err != nil {
		return "", err
	}

	if sl.activeRotations.Contains(pattern) {
		// Exact match
		return pattern, nil
	}

	ids := []string{}
	re, err := regexp.Compile(`.*` + pattern + `.*`)
	if err != nil {
		return "", err
	}
	sl.activeRotations.ForEach(func(id string) {
		if re.MatchString(id) {
			ids = append(ids, id)
		}
	})

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
	sl.activeRotations.Delete(r.RotationID)
	err = sl.Store.Index(KeyActiveRotations).Store(sl.activeRotations)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", r.RotationID)
	}

	logger.Infof("%s archived rotation %s.", sl.actingUser.Markdown(), r.Markdown())
	return nil
}

func (sl *sl) DebugDeleteRotation(rotationID string) error {
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
	sl.activeRotations.Delete(rotationID)
	err = sl.Store.Index(KeyActiveRotations).Store(sl.activeRotations)
	if err != nil {
		return errors.WithMessagef(err, "failed to store rotation %s", rotationID)
	}

	logger.Infof("%s deleted rotation %s.", sl.actingUser.Markdown(), rotationID)
	return nil
}

func (sl *sl) LoadRotation(rotationID string) (*Rotation, error) {
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

	return r, nil
}

func (sl *sl) MakeRotation(rotationName string) (*Rotation, error) {
	// id := ""
	// for i := 0; i < 5; i++ {
	// 	tryId := rotationName + "-" + model.NewId()[:7]
	// 	if len(sl.activeRotations) == 0 || sl.activeRotations[tryId] == "" {
	// 		id = tryId
	// 		break
	// 	}
	// }
	// if id == "" {
	// 	return nil, errors.New("Failed to generate unique rotation ID")
	// }

	// rotation := NewRotation(rotationName)
	// rotation.RotationID = id
	// return rotation, nil
	return nil, nil
}

func (sl *sl) ExpandRotation(r *Rotation) error {
	if len(r.users) != r.MattermostUserIDs.Len() {
		users, err := sl.LoadStoredUsers(r.MattermostUserIDs)
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
		user.LastServed[r.RotationID] = starting.Unix()

		user, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return added, err
		}

		if r.MattermostUserIDs == nil {
			r.MattermostUserIDs = types.NewSet()
		}
		r.MattermostUserIDs.Add(user.MattermostUserID)
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

		delete(user.LastServed, r.RotationID)
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
