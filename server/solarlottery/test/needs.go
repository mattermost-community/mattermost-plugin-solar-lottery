// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"

const (
	SkillServer  = "server"
	SkillWebapp  = "webapp"
	SkillMobile  = "mobile"
	SkillPlugins = "plugins"
)

func NeedMobile_L1_Min1() *store.Need { return store.NewNeed(SkillMobile, 1, 1) }
func NeedMobile_L1_Min2() *store.Need { return store.NewNeed(SkillMobile, 1, 2) }
func NeedMobile_L1_Min3() *store.Need { return store.NewNeed(SkillMobile, 1, 3) }

func NeedServer_L1_Min1() *store.Need { return store.NewNeed(SkillServer, 1, 1) }
func NeedServer_L1_Min3() *store.Need { return store.NewNeed(SkillServer, 1, 3) }
func NeedServer_L2_Min1() *store.Need { return store.NewNeed(SkillServer, 2, 1) }
func NeedServer_L2_Min2() *store.Need { return store.NewNeed(SkillServer, 2, 2) }

func NeedWebapp_L1_Min1() *store.Need { return store.NewNeed(SkillWebapp, 1, 1) }
func NeedWebapp_L2_Min1() *store.Need { return store.NewNeed(SkillWebapp, 2, 1) }
