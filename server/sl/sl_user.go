// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (sl *sl) ActingUser() (*User, error) {
	err := sl.Setup(withExpandedActingUser)
	if err != nil {
		return nil, err
	}
	return sl.actingUser, nil
}

func (sl *sl) LoadMattermostUserByUsername(username string) (*User, error) {
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

func (sl *sl) LoadUsers(mattermostUserIDs *types.IDSet) (*Users, error) {
	users, err := sl.loadStoredUsers(mattermostUserIDs)
	if err != nil {
		return nil, err
	}
	err = sl.expandUsers(users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (sl *sl) qualify(users *Users, skillLevels []SkillLevel) error {
	for _, skillLevel := range skillLevels {
		err := sl.AddKnownSkill(skillLevel.Skill)
		if err != nil {
			return err
		}
	}

	for _, user := range users.AsArray() {
		updated := []SkillLevel{}
		for _, skillLevel := range skillLevels {
			newSkill, newLevel := skillLevel.Skill, skillLevel.Level
			if !user.SkillLevels.Contains(newSkill) || Level(user.SkillLevels.Get(newSkill)) != newLevel {
				user.SkillLevels.Set(newSkill, int64(newLevel))
				updated = append(updated, skillLevel)
			}
		}
		if len(updated) == 0 {
			return nil
		}

		err := sl.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
		sl.Debugf("qualified %s for %s", user, updated)
	}
	return nil
}

func (sl *sl) disqualify(users *Users, skillNames []string) error {
	for _, user := range users.AsArray() {
		updated := []types.ID{}
		for _, skill := range skillNames {
			if user.SkillLevels.Contains(types.ID(skill)) {
				user.SkillLevels.Delete(types.ID(skill))
				updated = append(updated, types.ID(skill))
			}
		}
		if len(updated) == 0 {
			return nil
		}

		err := sl.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
		sl.Debugf("disqualified %s from %s", user, updated)
	}
	return nil
}

func (sl *sl) joinRotation(users *Users, r *Rotation, starting types.Time) (added *Users, err error) {
	added = NewUsers()

	for _, user := range users.AsArray() {
		if r.MattermostUserIDs.Contains(user.MattermostUserID) {
			sl.Debugf("%s is already in rotation %s.",
				added.MarkdownWithSkills(), r.Markdown())
			continue
		}

		// A new person may be given some slack - setting starting in the
		// future all but guarantees they won't be selected until then.
		user.LastServed.Set(r.RotationID, starting.Unix())

		err = sl.storeUserWelcomeNew(user)
		if err != nil {
			return added, err
		}

		if r.MattermostUserIDs == nil {
			r.MattermostUserIDs = types.NewIDSet()
		}
		r.MattermostUserIDs.Set(user.MattermostUserID)
		sl.dmUserWelcomeToRotation(user, r)
		added.Set(user)
	}
	return added, nil
}

func (sl *sl) leaveRotation(users *Users, r *Rotation) (*Users, error) {
	deleted := NewUsers()
	for _, user := range users.AsArray() {
		if !r.MattermostUserIDs.Contains(user.MattermostUserID) {
			sl.Debugf("%s is not found in rotation %s", user.Markdown(), r.Markdown())
			continue
		}

		user.LastServed.Delete(r.RotationID)
		err := sl.storeUserWelcomeNew(user)
		if err != nil {
			return nil, err
		}
		r.MattermostUserIDs.Delete(user.MattermostUserID)
		sl.dmUserLeftRotation(user, r)
		deleted.Set(user)
	}
	return deleted, nil
}

func (sl *sl) loadOrMakeUser(mattermostUserID types.ID) (*User, bool, error) {
	user, err := sl.loadUser(mattermostUserID)
	if err == kvstore.ErrNotFound {
		newUser := NewUser(mattermostUserID)
		newErr := sl.storeUserWelcomeNew(newUser)
		return newUser, true, newErr
	}
	if err != nil {
		return nil, false, err
	}
	user.loaded = true
	return user, false, nil
}

func (sl *sl) loadUser(mattermostUserID types.ID) (*User, error) {
	user := NewUser("")
	err := sl.Store.Entity(KeyUser).Load(mattermostUserID, user)
	if err != nil {
		return nil, err
	}
	user.loaded = true
	return user, nil
}

func (sl *sl) expandUser(user *User) error {
	if user == nil || user.MattermostUserID == "" {
		return errors.New("unreachable: expandUser: nil or no ID")
	}
	if !user.loaded {
		loaded, err := sl.loadUser(user.MattermostUserID)
		if err != nil {
			return err
		}
		*user = *loaded
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

func (sl *sl) expandUsers(users *Users) error {
	for _, user := range users.AsArray() {
		err := sl.expandUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

// storeUserWelcomeNew checks if the user being stored is new, and welcomes the user.
// note that it can be used inside of filters, so it must not use filters itself,
//  nor assume that any runtime values have been filled.
func (sl *sl) storeUserWelcomeNew(user *User) error {
	userIsNew := false
	if user.PluginVersion == "" {
		userIsNew = true
	}
	err := sl.storeUser(user)
	if err != nil {
		return err
	}
	if userIsNew {
		sl.dmUserWelcomeToSolarLottery(user)
	}
	return nil
}

func (sl *sl) storeUser(user *User) error {
	user.PluginVersion = sl.Config().PluginVersion
	err := sl.Store.Entity(KeyUser).Store(user.MattermostUserID, user)
	if err != nil {
		return err
	}
	return nil
}

func (sl *sl) loadStoredUsers(ids *types.IDSet) (*Users, error) {
	users := NewUsers()
	for _, id := range ids.IDs() {
		user, err := sl.loadUser(id)
		if err != nil {
			return NewUsers(), err
		}
		users.Set(user)
	}
	return users, nil
}

func (sl *sl) storeUsers(users *Users) error {
	for _, user := range users.AsArray() {
		err := sl.storeUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sl *sl) addUserUnavailable(user *User, u *Unavailable) error {
	user.AddUnavailable(u)
	err := sl.storeUser(user)
	if err != nil {
		return errors.Wrapf(err, "user: %s", user.Markdown())
	}
	return nil
}
