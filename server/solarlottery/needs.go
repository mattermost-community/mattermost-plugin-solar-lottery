// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func IsUserQualifiedForNeed(user *User, need *store.Need) bool {
	skillLevel, _ := user.SkillLevels[need.Skill]
	return skillLevel >= need.Level
}

func UsersQualifiedForNeed(users UserMap, need *store.Need) UserMap {
	qualified := UserMap{}
	for id, user := range users {
		if IsUserQualifiedForNeed(user, need) {
			qualified[id] = user
		}
	}
	return qualified
}

func UnmetNeeds(needs store.Needs, users UserMap) store.Needs {
	work := append(store.Needs{}, needs...)
	for i, need := range work {
		for _, user := range users {
			if IsUserQualifiedForNeed(user, need) {
				work[i].Min--
				work[i].Max--
			}
		}
	}

	var unmet store.Needs
	for _, need := range work {
		if need.Min > 0 {
			unmet = append(unmet, need)
		}
	}
	return unmet
}
