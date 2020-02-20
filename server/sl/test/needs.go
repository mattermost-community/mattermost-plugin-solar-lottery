// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"

const (
	SkillServer  = "server"
	SkillWebapp  = "webapp"
	SkillMobile  = "mobile"
	SkillPlugins = "plugins"
)

func NeedMobile_L1_Min1() *sl.Need { return sl.NewNeed(SkillMobile, 1, 1) }
func NeedMobile_L1_Min2() *sl.Need { return sl.NewNeed(SkillMobile, 1, 2) }
func NeedMobile_L1_Min3() *sl.Need { return sl.NewNeed(SkillMobile, 1, 3) }

func NeedServer_L1_Min1() *sl.Need { return sl.NewNeed(SkillServer, 1, 1) }
func NeedServer_L1_Min3() *sl.Need { return sl.NewNeed(SkillServer, 1, 3) }
func NeedServer_L2_Min1() *sl.Need { return sl.NewNeed(SkillServer, 2, 1) }
func NeedServer_L2_Min2() *sl.Need { return sl.NewNeed(SkillServer, 2, 2) }

func NeedWebapp_L1_Min1() *sl.Need { return sl.NewNeed(SkillWebapp, 1, 1) }
func NeedWebapp_L2_Min1() *sl.Need { return sl.NewNeed(SkillWebapp, 2, 1) }
