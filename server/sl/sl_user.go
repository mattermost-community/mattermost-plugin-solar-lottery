// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const eUser = "user_"

type UserService interface {
	LoadMattermostUsername(username string) (*User, error)
	AddToCalendar(mattermostUserIDs *types.IDIndex, u *Unavailable) (Users, error)
	ClearCalendar(mattermostUserIDs *types.IDIndex, interval types.Interval) (Users, error)
	Disqualify(mattermostUserIDs *types.IDIndex, skillName types.ID) (Users, error)
	Qualify(mattermostUserIDs *types.IDIndex, skillLevel SkillLevel) (Users, error)
	JoinRotation(mattermostUserIDs *types.IDIndex, rotationID types.ID, starting types.Time) (Users, error)
	LeaveRotation(mattermostUserIDs *types.IDIndex, rotationID types.ID) (Users, error)
	LoadUsers(mattermostUserIDs *types.IDIndex) (Users, error)
}

func (sl *sl) ActingUser() (*User, error) {
	err := sl.Setup(withExpandedActingUser)
	if err != nil {
		return nil, err
	}
	return sl.actingUser, nil
}

func (sl *sl) LoadMattermostUsername(username string) (*User, error) {
	if strings.HasPrefix(username, "@") {
		username = username[1:]
	}
	mmuser, err := sl.PluginAPI.GetMattermostUserByUsername(username)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to load %s", username)
	}
	user, _, err := sl.loadOrMakeUser(types.ID(mmuser.Id))
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to load User %s", username)
	}
	user.mattermostUser = mmuser
	return user, nil
}

func (sl *sl) LoadUsers(mattermostUserIDs *types.IDIndex) (Users, error) {
	users, err := sl.loadStoredUsers(mattermostUserIDs)
	if err != nil {
		return users, err
	}
	err = sl.expandUsers(users)
	return users, err
}

func (sl *sl) Qualify(mattermostUserIDs *types.IDIndex, skillLevel SkillLevel) (Users, error) {
	var users Users
	err := sl.Setup(
		pushLogger("Qualify", bot.LogContext{ctxSkillLevel: skillLevel}),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return users, err
	}
	defer sl.popLogger()

	err = sl.AddKnownSkill(skillLevel.Skill)
	if err != nil {
		return users, err
	}
	for _, user := range users.AsArray() {
		err = sl.updateUserSkill(user, skillLevel)
		if err != nil {
			return users, err
		}
	}

	sl.Infof("%s added skill %s to %s.",
		sl.actingUser.Markdown(), skillLevel, users.Markdown())
	return users, nil
}

func (sl *sl) Disqualify(mattermostUserIDs *types.IDIndex, skillName types.ID) (Users, error) {
	var users Users
	err := sl.Setup(
		pushLogger("Disqualify", bot.LogContext{ctxSkill: skillName}),
		withValidSkillName(skillName),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return users, err
	}
	defer sl.popLogger()

	for _, user := range users.AsArray() {
		err = sl.updateUserSkill(user, NewSkillLevel(skillName, 0))
		if err != nil {
			return users, err
		}
	}

	sl.Infof("%s removed skill %s from %s.",
		sl.actingUser.Markdown(), skillName, users.Markdown())
	return users, nil
}

func (sl *sl) JoinRotation(mattermostUserIDs *types.IDIndex, rotationID types.ID, starting types.Time) (Users, error) {
	var users Users
	err := sl.Setup(
		pushLogger("JoinRotation", bot.LogContext{ctxStarting: starting}),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return NewUsers(), err
	}
	defer sl.popLogger()

	added := NewUsers()

	r, err := sl.UpdateRotation(rotationID, func(r *Rotation) error {
		for _, user := range users.AsArray() {
			if r.MattermostUserIDs.Contains(user.MattermostUserID) {
				sl.Debugf("%s is already in rotation %s.",
					added.MarkdownWithSkills(), r.Markdown())
				continue
			}

			// A new person may be given some slack - setting starting in the
			// future all but guarantees they won't be selected until then.
			user.LastServed.Set(r.RotationID, starting.Unix())

			user, err = sl.storeUserWelcomeNew(user)
			if err != nil {
				return err
			}

			if r.MattermostUserIDs == nil {
				r.MattermostUserIDs = types.NewIDIndex()
			}
			r.MattermostUserIDs.Set(user.MattermostUserID)
			sl.messageWelcomeToRotation(user, r)
			added.Set(user)
		}
		return nil
	})

	sl.Infof("%s added %s to %s.",
		sl.actingUser.Markdown(), added.MarkdownWithSkills(), r.Markdown())
	return added, nil
}

func (sl *sl) LeaveRotation(mattermostUserIDs *types.IDIndex, rotationID types.ID) (Users, error) {
	var r *Rotation
	var users Users
	err := sl.Setup(
		pushLogger("LeaveRotation", nil),
		withRotation(rotationID, &r),
		withExpandedUsers(mattermostUserIDs, &users),
	)
	if err != nil {
		return users, err
	}
	defer sl.popLogger()

	deleted := NewUsers()
	for _, user := range users.AsArray() {
		if !r.MattermostUserIDs.Contains(user.MattermostUserID) {
			sl.Debugf("%s is not found in rotation %s", user.Markdown(), r.Markdown())
			continue
		}

		user.LastServed.Delete(r.RotationID)
		_, err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return deleted, err
		}
		r.MattermostUserIDs.Delete(user.MattermostUserID)
		sl.messageLeftRotation(user, r)
		deleted.Set(user)
	}

	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return deleted, err
	}

	sl.Infof("%s removed from %s.", deleted.Markdown(), r.Markdown())
	return deleted, nil
}

func (sl *sl) loadOrMakeUser(mattermostUserID types.ID) (*User, bool, error) {
	var user User
	err := sl.Store.Entity(KeyUser).Load(mattermostUserID, &user)
	if err == kvstore.ErrNotFound {
		newUser, newErr := sl.storeUserWelcomeNew(NewUser(mattermostUserID))
		return newUser, true, newErr
	}
	if err != nil {
		return nil, false, err
	}
	user.loaded = true
	return &user, false, nil
}

func (sl *sl) expandUser(user *User) error {
	if user == nil || user.MattermostUserID == "" {
		return errors.New("unreachable: expandUser: nil or no ID")
	}
	if !user.loaded {
		err := sl.Store.Entity(KeyUser).Load(user.MattermostUserID, user)
		if err != nil {
			return err
		}
		user.loaded = true
	}
	if user.mattermostUser == nil {
		mattermostUser, err := sl.PluginAPI.GetMattermostUser(string(user.MattermostUserID))
		if err != nil {
			return err
		}
		user.mattermostUser = mattermostUser

		loc, err := time.LoadLocation(mattermostUser.GetPreferredTimezone())
		if err != nil {
			return err
		}
		user.location = loc
	}
	return nil
}

// storeUserWelcomeNew checks if the user being stored is new, and welcomes the user.
// note that it can be used inside of filters, so it must not use filters itself,
//  nor assume that any runtime values have been filled.
func (sl *sl) storeUserWelcomeNew(user *User) (*User, error) {
	if user.PluginVersion == "" {
		sl.messageWelcomeNewUser(user)
	}
	err := sl.storeUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (sl *sl) storeUser(user *User) error {
	user.PluginVersion = sl.Config().PluginVersion
	err := sl.Store.Entity(KeyUser).Store(user.MattermostUserID, user)
	if err != nil {
		return err
	}
	return nil
}

func (sl *sl) updateUserSkill(user *User, skillLevel SkillLevel) error {
	s, l := skillLevel.Skill, skillLevel.Level
	if user.SkillLevels.Contains(s) && Level(user.SkillLevels.Get(s)) == l {
		// nothing to do
		sl.Debugf("nothing to do for user %s, already is %s", user.Markdown(), skillLevel)
		return nil
	}

	if l == 0 {
		user.SkillLevels.Delete(s)
	} else {
		user.SkillLevels.Set(s, int64(l))
	}
	user, err := sl.storeUserWelcomeNew(user)
	if err != nil {
		return err
	}
	sl.Debugf("%s updated to %s", user.Markdown(), skillLevel)
	return nil
}

func (sl *sl) loadStoredUsers(ids *types.IDIndex) (Users, error) {
	users := NewUsers()
	for _, id := range ids.IDs() {
		var user User
		err := sl.Store.Entity(KeyUser).Load(id, &user)
		if err != nil {
			return NewUsers(), err
		}
		users.Set(&user)
	}
	return users, nil
}
