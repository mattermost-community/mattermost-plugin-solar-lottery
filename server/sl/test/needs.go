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

func MobileL1() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 1) }
func MobileL2() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 2) }
func MobileL3() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 3) }
func MobileL4() sl.SkillLevel { return sl.NewSkillLevel(Mobile, 4) }

func ServerL1() sl.SkillLevel { return sl.NewSkillLevel(Server, 1) }
func ServerL2() sl.SkillLevel { return sl.NewSkillLevel(Server, 2) }
func ServerL3() sl.SkillLevel { return sl.NewSkillLevel(Server, 3) }
func ServerL4() sl.SkillLevel { return sl.NewSkillLevel(Server, 4) }

func WebappL1() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 1) }
func WebappL2() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 2) }
func WebappL3() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 3) }
func WebappL4() sl.SkillLevel { return sl.NewSkillLevel(Webapp, 4) }

func C1Any() sl.Need { return sl.NewNeed(1, sl.AnySkillLevel) }
func C2Any() sl.Need { return sl.NewNeed(1, sl.AnySkillLevel) }
func C3Any() sl.Need { return sl.NewNeed(1, sl.AnySkillLevel) }

func C1MobileL1() sl.Need { return sl.NewNeed(1, MobileL1()) }
func C1MobileL2() sl.Need { return sl.NewNeed(1, MobileL2()) }
func C1MobileL3() sl.Need { return sl.NewNeed(1, MobileL3()) }
func C1MobileL4() sl.Need { return sl.NewNeed(1, MobileL4()) }
func C2MobileL1() sl.Need { return sl.NewNeed(2, MobileL1()) }
func C2MobileL2() sl.Need { return sl.NewNeed(2, MobileL2()) }
func C2MobileL3() sl.Need { return sl.NewNeed(2, MobileL3()) }
func C2MobileL4() sl.Need { return sl.NewNeed(2, MobileL4()) }

func C1ServerL1() sl.Need { return sl.NewNeed(1, ServerL1()) }
func C1ServerL2() sl.Need { return sl.NewNeed(1, ServerL2()) }
func C1ServerL3() sl.Need { return sl.NewNeed(1, ServerL3()) }
func C1ServerL4() sl.Need { return sl.NewNeed(1, ServerL4()) }
func C2ServerL1() sl.Need { return sl.NewNeed(2, ServerL1()) }
func C2ServerL2() sl.Need { return sl.NewNeed(2, ServerL2()) }
func C2ServerL3() sl.Need { return sl.NewNeed(2, ServerL3()) }
func C2ServerL4() sl.Need { return sl.NewNeed(2, ServerL4()) }
func C3_ServerL1() sl.Need { return sl.NewNeed(3, ServerL1()) }

func C1WebappL1() sl.Need { return sl.NewNeed(1, WebappL1()) }
func C1WebappL2() sl.Need { return sl.NewNeed(1, WebappL2()) }
func C1WebappL3() sl.Need { return sl.NewNeed(1, WebappL3()) }
func C1WebappL4() sl.Need { return sl.NewNeed(1, WebappL4()) }
func C2WebappL1() sl.Need { return sl.NewNeed(2, WebappL1()) }
func C2WebappL2() sl.Need { return sl.NewNeed(2, WebappL2()) }
func C2WebappL3() sl.Need { return sl.NewNeed(2, WebappL3()) }
func C2WebappL4() sl.Need { return sl.NewNeed(2, WebappL4()) }
