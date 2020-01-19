// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"

const (
	testSkillServer  = "server"
	testSkillWebapp  = "webapp"
	testSkillMobile  = "mobile"
	testSkillPlugins = "plugins"
)

func testNeedMobile_L1_Min1() *store.Need { return store.NewNeed(testSkillMobile, 1, 1) }
func testNeedMobile_L1_Min2() *store.Need { return store.NewNeed(testSkillMobile, 1, 2) }
func testNeedMobile_L1_Min3() *store.Need { return store.NewNeed(testSkillMobile, 1, 3) }

func testNeedServer_L1_Min1() *store.Need { return store.NewNeed(testSkillServer, 1, 1) }
func testNeedServer_L1_Min3() *store.Need { return store.NewNeed(testSkillServer, 1, 3) }
func testNeedServer_L2_Min1() *store.Need { return store.NewNeed(testSkillServer, 2, 1) }
func testNeedServer_L2_Min2() *store.Need { return store.NewNeed(testSkillServer, 2, 2) }

func testNeedWebapp_L1_Min1() *store.Need { return store.NewNeed(testSkillWebapp, 1, 1) }
func testNeedWebapp_L2_Min1() *store.Need { return store.NewNeed(testSkillWebapp, 2, 1) }
