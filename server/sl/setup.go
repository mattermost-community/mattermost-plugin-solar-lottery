// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type filterf func(*sl) error

func (sl *sl) Setup(filters ...filterf) error {
	for _, filter := range filters {
		err := filter(sl)
		if err != nil {
			return err
		}
	}
	return nil
}

func withRotation(rotationID types.ID, ref **Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		sl.Logger = sl.Logger.With(bot.LogContext{ctxRotationID: rotationID})

		r, err := sl.loadRotation(rotationID)
		if err != nil {
			return err
		}
		*ref = r
		return nil
	}
}

func withExpandedRotation(rotationID types.ID, ref **Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		err := withRotation(rotationID, ref)(sl)
		if err != nil {
			return err
		}
		return sl.expandRotation(*ref)
	}
}

func withRotationIsActive(r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		if r.IsArchived {
			return errors.Errorf("rotation %s is archived", r.Markdown())
		}
		return nil
	}
}

func withUser(mattermostUserID types.ID, userref **User) func(sl *sl) error {
	return func(sl *sl) error {
		loadedUser, _, err := sl.loadOrMakeUser(mattermostUserID)
		if err != nil {
			return err
		}
		*userref = loadedUser
		return nil
	}
}

func withExpandedUser(mattermostUserID types.ID, ref **User) func(sl *sl) error {
	return func(sl *sl) error {
		err := withUser(mattermostUserID, ref)(sl)
		if err != nil {
			return err
		}
		return sl.expandUser(*ref)
	}
}

func withActingUser(sl *sl) error {
	sl.Logger = sl.Logger.With(bot.LogContext{ctxActingUserID: sl.actingMattermostUserID})
	return withUser(sl.actingMattermostUserID, &sl.actingUser)(sl)
}

func withExpandedActingUser(sl *sl) error {
	err := withExpandedUser(sl.actingMattermostUserID, &sl.actingUser)(sl)
	if err != nil {
		return err
	}
	sl.Logger = sl.Logger.With(bot.LogContext{ctxActingUsername: sl.actingUser.MattermostUsername()})
	return nil
}

func withUsers(mattermostUserIDs *types.IDIndex, ref *Users) func(sl *sl) error {
	return func(sl *sl) error {
		users, err := sl.loadStoredUsers(mattermostUserIDs)
		if err != nil {
			return err
		}
		*ref = users
		sl.Logger = sl.Logger.With(bot.LogContext{ctxUserIDs: users.IDs()})
		return nil
	}
}

func withExpandedUsers(mattermostUserIDs *types.IDIndex, ref *Users) func(sl *sl) error {
	return func(sl *sl) error {
		err := withUsers(mattermostUserIDs, ref)(sl)
		if err != nil {
			return err
		}
		users := *ref
		if users.IsEmpty() {
			return nil
		}
		err = sl.expandUsers(users)
		sl.Logger = sl.Logger.With(bot.LogContext{ctxUsernames: users.String()})
		return nil
	}
}

func withKnownSkills(ref **types.IDIndex) func(*sl) error {
	return func(sl *sl) error {
		skills, err := sl.Store.IDIndex(KeyKnownSkills).Load()
		if err == kvstore.ErrNotFound {
			*ref = types.NewIDIndex()
			return nil
		}
		if err != nil {
			return err
		}
		*ref = skills
		return nil
	}
}

func withValidSkillName(skillName types.ID) func(sl *sl) error {
	return func(sl *sl) error {
		var knownSkills *types.IDIndex
		err := sl.Setup(withKnownSkills(&knownSkills))
		if err != nil {
			return err
		}
		if !knownSkills.Contains(skillName) {
			return errors.Errorf("skill %s is not found", skillName)
		}
		return nil
	}
}

func withActiveRotations(ref **types.IDIndex) func(sl *sl) error {
	return func(sl *sl) error {
		var activeRotations *types.IDIndex
		activeRotations, err := sl.Store.IDIndex(KeyActiveRotations).Load()
		if err == kvstore.ErrNotFound {
			*ref = types.NewIDIndex()
			return nil
		}
		if err != nil {
			return err
		}
		*ref = activeRotations
		return nil
	}
}

func pushLogger(apiName string, logContext bot.LogContext) func(*sl) error {
	return func(sl *sl) error {
		err := withExpandedActingUser(sl)
		if err != nil {
			return err
		}

		logger := sl.Logger
		logger = logger.With(logContext)
		logger = logger.With(bot.LogContext{
			ctxActingUsername: sl.actingUser.MattermostUsername(),
			ctxAPI:            apiName,
		})

		if sl.loggers == nil {
			sl.loggers = []bot.Logger{logger}
		} else {
			sl.loggers = append(sl.loggers, logger)
		}
		sl.Logger = logger
		return nil
	}
}

func (sl *sl) popLogger() {
	l := len(sl.loggers)
	if l == 0 {
		return
	}
	sl.Logger = sl.loggers[l-1]
	sl.loggers = sl.loggers[:l-1]
}
