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

func withLoadActiveRotations(activeRotations *types.IDSet) func(sl *sl) error {
	return func(sl *sl) error {
		loaded, err := sl.Store.IDIndex(KeyActiveRotations).Load()
		if err == kvstore.ErrNotFound {
			*activeRotations = *types.NewIDSet()
			return nil
		}
		if err != nil {
			return err
		}
		*activeRotations = *loaded
		return nil
	}
}

func withLoadRotation(rotationID *types.ID, r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		sl.Logger = sl.Logger.With(bot.LogContext{ctxRotationID: rotationID})

		loaded, err := sl.loadRotation(*rotationID)
		if err != nil {
			return err
		}
		*r = *loaded
		return nil
	}
}

func withExpandRotationUsers(r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		return sl.expandRotationUsers(r)
	}
}

func withExpandRotationTasks(r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		return sl.expandRotationTasks(r)
	}
}

func withExpandedRotation(rotationID types.ID, r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		return sl.Setup(
			withLoadRotation(&rotationID, r),
			withExpandRotationUsers(r),
			withExpandRotationTasks(r),
		)
	}
}

func withLoadOrMakeUser(mattermostUserID *types.ID, user *User) func(sl *sl) error {
	return func(sl *sl) error {
		loadedUser, _, err := sl.loadOrMakeUser(*mattermostUserID)
		if err != nil {
			return err
		}
		*user = *loadedUser
		return nil
	}
}

func withExpandedUser(mattermostUserID types.ID, user *User) func(sl *sl) error {
	return func(sl *sl) error {
		err := withLoadOrMakeUser(&mattermostUserID, user)(sl)
		if err != nil {
			return err
		}
		return sl.expandUser(user)
	}
}

func withExpandedActingUser(sl *sl) error {
	sl.actingUser = NewUser("")
	err := sl.Setup(withExpandedUser(sl.actingMattermostUserID, sl.actingUser))
	if err != nil {
		return err
	}
	sl.Logger = sl.Logger.With(bot.LogContext{ctxActingUsername: sl.actingUser.MattermostUsername()})
	return nil
}

func withLoadUsers(mattermostUserIDs **types.IDSet, users *Users) func(sl *sl) error {
	return func(sl *sl) error {
		loaded, err := sl.loadStoredUsers(*mattermostUserIDs)
		if err != nil {
			return err
		}
		*users = *loaded
		sl.Logger = sl.Logger.With(bot.LogContext{ctxUserIDs: users.IDs()})
		return nil
	}
}

func withExpandedUsers(mattermostUserIDs **types.IDSet, users *Users) func(sl *sl) error {
	return func(sl *sl) error {
		err := withLoadUsers(mattermostUserIDs, users)(sl)
		if err != nil {
			return err
		}
		if users.IsEmpty() {
			return nil
		}
		err = sl.expandUsers(users)
		if err != nil {
			return err
		}
		sl.Logger = sl.Logger.With(bot.LogContext{ctxUsernames: users.String()})
		return nil
	}
}

func withLoadKnownSkills(knownSkills *types.IDSet) func(*sl) error {
	return func(sl *sl) error {
		skills, err := sl.Store.IDIndex(KeyKnownSkills).Load()
		if err == kvstore.ErrNotFound {
			*knownSkills = *types.NewIDSet()
			return nil
		}
		if err != nil {
			return err
		}
		*knownSkills = *skills
		return nil
	}
}

func withValidSkillName(skillName *types.ID) func(sl *sl) error {
	return func(sl *sl) error {
		knownSkills := *types.NewIDSet()
		err := sl.Setup(withLoadKnownSkills(&knownSkills))
		if err != nil {
			return err
		}
		if !knownSkills.Contains(*skillName) {
			return errors.Errorf("skill %s is not found", skillName)
		}
		return nil
	}
}

func withLoadTask(taskID *types.ID, task *Task) func(sl *sl) error {
	return func(sl *sl) error {
		sl.Logger = sl.Logger.With(bot.LogContext{ctxTaskID: *taskID})

		loaded, err := sl.loadTask(*taskID)
		if err != nil {
			return err
		}
		*task = *loaded
		return nil
	}
}

func pushAPILogger(apiName string, in interface{}) func(*sl) error {
	return func(sl *sl) error {
		err := withExpandedActingUser(sl)
		if err != nil {
			return err
		}

		logger := sl.Logger
		logger = logger.With(bot.LogContext{
			ctxActingUsername: sl.actingUser.MattermostUsername(),
			ctxAPI:            apiName,
			ctxInput:          in,
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
