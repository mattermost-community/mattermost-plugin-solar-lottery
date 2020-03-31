// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package plugin

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const commandCleanInit = `
/lotto debug-clean

/lotto user qualify -s ABC-FS1 -l intermediate -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/lotto user qualify -s DEF-PERF -l intermediate -u @deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott
/lotto user qualify -s GHIJ-FS2 -l intermediate -u @helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson
/lotto user qualify -s KL-SRE -l intermediate -u @karen.austin,@karen.martin,@kathryn.mills,@laura.wagner
/lotto user qualify -s MN-FS3 -l intermediate -u @margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts
/lotto user qualify -s R-MOBILE -l intermediate -u @ralph.watson,@raymond.austin,@raymond.fisher,@raymond.fox

/lotto user qualify -s lead -l intermediate -u @deborah.freeman,@karen.austin,@ralph.watson
/lotto user qualify -s lead -l advanced -u @aaron.medina,@helen.hunter,@margaret.morgan

/lotto user qualify -s server -l beginner -u @aaron.peterson,@albert.torres,@helen.hunter,@johnny.hansen,@mark.rodriguez,@raymond.fox
/lotto user qualify -s server -l intermediate -u @alice.johnston,@eugene.rodriguez,@janice.armstrong,@jerry.ramos,@jonathan.watson,@kathryn.mills,@margaret.morgan,@matthew.mendoza,@mildred.barnes
/lotto user qualify -s server -l advanced -u @aaron.medina,@deborah.freeman,@diana.wells,@emily.meyer,@nancy.roberts
/lotto user qualify -s server -l expert -u @douglas.daniels

/lotto user qualify -s webapp -l beginner -u @emily.meyer,@eugene.rodriguez,@helen.hunter,@janice.armstrong,@johnny.hansen,@karen.martin,@mark.rodriguez,@mildred.barnes,@raymond.austin,@raymond.fisher,@raymond.fox
/lotto user qualify -s webapp -l intermediate -u @aaron.peterson,@aaron.ward,@deborah.freeman,@jeremy.williamson,@jerry.ramos,@nancy.roberts
/lotto user qualify -s webapp -l advanced -u @aaron.medina,@frances.elliott,@matthew.mendoza

/lotto user qualify -s mobile -l beginner -u @johnny.hansen,@aaron.peterson,@nancy.roberts
/lotto user qualify -s mobile -l intermediate -u @raymond.austin,@raymond.fisher,@raymond.fox
/lotto user qualify -s mobile -l advanced -u @ralph.watson,@frances.elliott

/lotto user qualify -s sre -l beginner -u @karen.martin
/lotto user qualify -s sre -l intermediate -u @laura.wagner
/lotto user qualify -s sre -l advanced -u @karen.austin,@kathryn.mills

/lotto user qualify -s build -l advanced -u @karen.martin
`

const commandIcebreakerInit = `
/lotto rotation add -r icebreaker --period w --start 2019-12-17 --grace 3 --size 2

/lotto rotation join -r icebreaker -s 2019-12-17 -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts,@ralph.watson,@raymond.austin,@raymond.fisher,@raymond.fox,@sysadmin
/lotto rotation show -r icebreaker
`

const commandIcebreakerForecast1 = `
/lotto rotation guess -r icebreaker -s 0 -n 20
`

const commandIcebreakerForecast2 = `
/lotto rotation forecast -r icebreaker -s 0 -n 20 --sample 100
`

const commandIcebreakerForecast = commandIcebreakerForecast1 + "\n" + commandIcebreakerForecast2

const commandIcebreakerAutopilot = `
/lotto rotation autopilot -r icebreaker --notify 3 --fill-before 5
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-11
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-12
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-13
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-17
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-21
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-25
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-29
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-30
`

const commandIcebreakerUser = `
/lotto user unavailable -s 2020-02-01 -e 2020-03-18 -u @helen.hunter
/lotto user forecast -n 20 -r ice -u @helen.hunter --sample 50
`

const commandSETInit = `
/lotto rotation add  -r SET --period m --start 2020-01-16 --grace 1 --size 5

/lotto rotation need -r SET --skill webapp --level beginner --min 2
/lotto rotation need -r SET --skill server --level beginner --min 2
/lotto rotation need -r SET --skill mobile --level beginner --min 1
/lotto rotation need -r SET --skill lead --level beginner --min 1 --max 1

/lotto rotation need -r SET --skill ABC-FS1 --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill GHIJ-FS2 --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill MN-FS3 --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill DEF-PERF --level beginner --min -1 --max 1
/lotto rotation need -r SET --skill KL-SRE --level beginner --min -1 --max 1
/lotto rotation need -r SET --skill R-MOBILE --level beginner --min -1 --max 1

/lotto rotation join -r SET -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts,@ralph.watson,@raymond.austin,@raymond.fisher,@raymond.fox

/lotto rotation show -r SET
`

const commandSETForecast1 = `
/lotto log --level debug 
/lotto rotation guess -r SET -s 0 -n 12
/lotto log --level info
`

const commandSETForecast2 = `
/lotto rotation forecast -r SET -s 0 -n 12 --sample 50
`

var constCommandAll = strings.Join([]string{
	commandCleanInit,
	commandIcebreakerInit,
	commandIcebreakerForecast,
	commandIcebreakerAutopilot,
	commandIcebreakerUser,
	commandSETInit,
	commandSETForecast,
}, "\n")

const commandSETForecast = commandSETForecast1 + "\n" + commandSETForecast2

func (p *Plugin) executeDemoCommand(c *plugin.Context, args *model.CommandArgs) bool {
	if c == nil || args == nil {
		return false
	}
	split := strings.Fields(args.Command)
	if len(split) != 3 {
		return false
	}
	if split[0] != "/"+constants.CommandTrigger || split[1] != "demo" {
		return false
	}

	run := ""
	switch split[2] {
	case "all":
		run = constCommandAll
	case "1", "clean-init":
		run = commandCleanInit
	case "2", "icebreaker-init":
		run = commandIcebreakerInit
	case "3", "icebreaker-forecast1":
		run = commandIcebreakerForecast1
	case "4", "icebreaker-forecast2":
		run = commandIcebreakerForecast2
	case "icebreaker-forecast":
		run = commandIcebreakerForecast
	case "5", "icebreaker-autopilot":
		run = commandIcebreakerAutopilot
	case "6", "icebreaker-user":
		run = commandIcebreakerUser
	case "7", "SET-init":
		run = commandSETInit
	case "8", "SET-forecast1":
		run = commandSETForecast1
	case "9", "SET-forecast2":
		run = commandSETForecast2
	case "SET-forecast":
		run = commandSETForecast
	}

	p.runCommands(c, args, run)
	return true
}

func (p *Plugin) runCommands(c *plugin.Context, args *model.CommandArgs, in string) {
	args = &(*args)
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		args.Command = line
		p.ExecuteCommand(c, args)
	}
}
