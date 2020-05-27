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

// key is always a constant (for the time being?) so no need for a double-pointer
func withLoadIDIndex(key string, idx *types.IDSet) func(sl *sl) error {
	return func(sl *sl) error {
		loaded, err := sl.Store.IDIndex(key).Load()
		if err == kvstore.ErrNotFound {
			idx.From(&types.NewIDSet().ValueSet)
			return nil
		}
		if err != nil {
			return err
		}
		idx.From(&loaded.ValueSet)
		return nil
	}
}

func withLoadRotation(idref *types.ID, r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		sl.Logger = sl.Logger.With(bot.LogContext{ctxRotationID: *idref})

		loaded, err := sl.loadRotation(*idref)
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

func withExpandedRotation(idref *types.ID, r *Rotation) func(sl *sl) error {
	return func(sl *sl) error {
		return sl.Setup(
			withLoadRotation(idref, r),
			withExpandRotationUsers(r),
			withExpandRotationTasks(r),
		)
	}
}

func withExpandedActingUser(sl *sl) error {
	user, _, err := sl.loadOrMakeUser(sl.actingMattermostUserID)
	if err != nil {
		return err
	}
	err = sl.expandUser(user)
	if err != nil {
		return err
	}
	sl.actingUser = user
	sl.Logger = sl.Logger.With(bot.LogContext{ctxActingUsername: user.MattermostUsername()})
	return nil
}

func withExpandedUsers(idsref **types.IDSet, users *Users) func(sl *sl) error {
	return func(sl *sl) error {
		loaded, err := sl.LoadUsers(*idsref)
		if err != nil {
			return err
		}
		users.From(&loaded.ValueSet)
		return nil
	}
}

func withValidSkillName(skillName *types.ID) func(sl *sl) error {
	return func(sl *sl) error {
		knownSkills := types.NewIDSet()
		err := sl.Setup(withLoadIDIndex(KeyKnownSkills, knownSkills))
		if err != nil {
			return err
		}
		if !knownSkills.Contains(*skillName) {
			return errors.Errorf("skill %s is not found", skillName)
		}
		return nil
	}
}

func withValidSkillNames(skillNames ...string) func(sl *sl) error {
	return func(sl *sl) error {
		for _, skill := range skillNames {
			id := types.ID(skill)
			err := withValidSkillName(&id)(sl)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func withExpandedTask(idref *types.ID, task *Task) func(sl *sl) error {
	return func(sl *sl) error {
		loaded, err := sl.LoadTask(*idref)
		if err != nil {
			return err
		}
		*task = *loaded
		sl.Logger = sl.Logger.With(bot.LogContext{ctxTaskID: task.TaskID})
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
