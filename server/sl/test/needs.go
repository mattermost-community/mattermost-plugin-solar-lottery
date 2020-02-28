// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package test

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"

const (
	Server  = "server"
	Webapp  = "webapp"
	Mobile  = "mobile"
	Plugins = "plugins"
)

func Mobile_L1() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 1) }
func Mobile_L2() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 2) }
func Mobile_L3() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 3) }
func Mobile_L4() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 4) }

func Server_L1() sl.SkillLevel { return sl.NewSkillLevel(Server, 1) }
func Server_L2() sl.SkillLevel { return sl.NewSkillLevel(Server, 2) }
func Server_L3() sl.SkillLevel { return sl.NewSkillLevel(Server, 3) }
func Server_L4() sl.SkillLevel { return sl.NewSkillLevel(Server, 4) }

func Webapp_L1() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 1) }
func Webapp_L2() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 2) }
func Webapp_L3() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 3) }
func Webapp_L4() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 4) }

func C1_Mobile_L1() *sl.Need { return sl.NewNeed(1, Mobile_L1()) }
func C1_Mobile_L2() *sl.Need { return sl.NewNeed(1, Mobile_L2()) }
func C1_Mobile_L3() *sl.Need { return sl.NewNeed(1, Mobile_L3()) }
func C1_Mobile_L4() *sl.Need { return sl.NewNeed(1, Mobile_L4()) }
func C2_Mobile_L1() *sl.Need { return sl.NewNeed(2, Mobile_L1()) }
func C2_Mobile_L2() *sl.Need { return sl.NewNeed(2, Mobile_L2()) }
func C2_Mobile_L3() *sl.Need { return sl.NewNeed(2, Mobile_L3()) }
func C2_Mobile_L4() *sl.Need { return sl.NewNeed(2, Mobile_L4()) }

func C1_Server_L1() *sl.Need { return sl.NewNeed(1, Server_L1()) }
func C1_Server_L2() *sl.Need { return sl.NewNeed(1, Server_L2()) }
func C1_Server_L3() *sl.Need { return sl.NewNeed(1, Server_L3()) }
func C1_Server_L4() *sl.Need { return sl.NewNeed(1, Server_L4()) }
func C2_Server_L1() *sl.Need { return sl.NewNeed(2, Server_L1()) }
func C2_Server_L2() *sl.Need { return sl.NewNeed(2, Server_L2()) }
func C2_Server_L3() *sl.Need { return sl.NewNeed(2, Server_L3()) }
func C2_Server_L4() *sl.Need { return sl.NewNeed(2, Server_L4()) }

func C1_Webapp_L1() *sl.Need { return sl.NewNeed(1, Webapp_L1()) }
func C1_Webapp_L2() *sl.Need { return sl.NewNeed(1, Webapp_L2()) }
func C1_Webapp_L3() *sl.Need { return sl.NewNeed(1, Webapp_L3()) }
func C1_Webapp_L4() *sl.Need { return sl.NewNeed(1, Webapp_L4()) }
func C2_Webapp_L1() *sl.Need { return sl.NewNeed(2, Webapp_L1()) }
func C2_Webapp_L2() *sl.Need { return sl.NewNeed(2, Webapp_L2()) }
func C2_Webapp_L3() *sl.Need { return sl.NewNeed(2, Webapp_L3()) }
func C2_Webapp_L4() *sl.Need { return sl.NewNeed(2, Webapp_L4()) }
