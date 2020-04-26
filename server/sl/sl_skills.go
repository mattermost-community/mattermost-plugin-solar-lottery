// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

func (sl *sl) ListKnownSkills() (*types.IDSet, error) {
	knownSkills := types.NewIDSet()
	err := sl.Setup(withLoadIDIndex(KeyKnownSkills, knownSkills))
	if err != nil {
		return nil, err
	}
	return knownSkills, nil
}

func (sl *sl) AddKnownSkill(skillName types.ID) error {
	err := sl.Setup(pushAPILogger("AddKnownSkill", skillName))
	if err != nil {
		return err
	}
	defer sl.popLogger()

	if skillName == AnySkill {
		return errors.Errorf("%s is reserved", skillName)
	}

	added, err := sl.Store.IDIndex(KeyKnownSkills).Set(skillName)
	if err != nil {
		return err
	}
	if added {
		sl.Infof("%s added known skill %s.", sl.actingUser.Markdown(), skillName)
	}
	return nil
}

func (sl *sl) DeleteKnownSkill(skillName types.ID) error {
	err := sl.Setup(pushAPILogger("DeleteKnownSkill", skillName))
	if err != nil {
		return err
	}
	defer sl.popLogger()

	err = sl.Store.IDIndex(KeyKnownSkills).Delete(skillName)
	if err != nil {
		return err
	}

	sl.Infof("%s deleted skill %s.", sl.actingUser.Markdown(), skillName)
	return nil
}
