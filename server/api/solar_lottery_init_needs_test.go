// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"

const (
	testSkillServer = "server"
	testSkillWebapp = "webapp"
	testSkillMobile = "mobile"
)

var (
	testNeedMobile_L1_Min1 = *store.NewNeed(testSkillMobile, 1, 1)
	testNeedMobile_L1_Min2 = *store.NewNeed(testSkillMobile, 1, 2)

	testNeedServer_L1_Min1 = *store.NewNeed(testSkillServer, 1, 1)
	testNeedServer_L1_Min3 = *store.NewNeed(testSkillServer, 1, 3)
	testNeedServer_L2_Min2 = *store.NewNeed(testSkillServer, 2, 2)

	testNeedWebapp_L1_Min1 = *store.NewNeed(testSkillWebapp, 1, 1)
	testNeedWebapp_L2_Min1 = *store.NewNeed(testSkillWebapp, 2, 1)
)
