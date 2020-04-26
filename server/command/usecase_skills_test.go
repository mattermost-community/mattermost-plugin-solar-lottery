// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"
)

func TestUseCaseSkillsSimple(t *testing.T) {
	multiseedBench(t, "TEST", 1, 3, 50, 0, `
		/lotto rotation new TEST --beginning=2019-01-16 --period=monthly
		/lotto rotation set task TEST --grace=300h 
		
		/lotto rotation set require TEST -s any --count 3
		/lotto rotation set require TEST -s S1 --count 2
		/lotto rotation set require TEST -s S2 --count 2
		
		/lotto user qualify @u01 @u02 @u03 -s TEAM1,S1
		/lotto user qualify @u04 @u05      -s TEAM1,S2
		/lotto user qualify @u06 @u07 @u08 -s TEAM2,S2
		/lotto user qualify @u09 @u10      -s TEAM2,S1
		/lotto user qualify @u11 @u12 @u13 -s TEAM3,S1
		/lotto user qualify @u14 @u15      -s TEAM4,S2
		/lotto user qualify @u16 @u17 @u18 -s TEAM4,S2
		/lotto user qualify @u19 @u20      -s TEAM4,S1
		
		/lotto user join TEST @u01 @u02 @u03 @u04 @u05 @u06 @u07 @u08 @u09 @u10 
		/lotto user join TEST @u11 @u12 @u13 @u14 @u15 @u16 @u17 @u18 @u19 @u20 
		`,
		nil,
		map[int][]string{
			13: []string{"@u14 (S2-◉, TEAM4-◉)", "@u20 (S1-◉, TEAM4-◉)"},
			15: []string{"@u04 (S2-◉, TEAM1-◉)"},
			16: []string{"@u19 (S1-◉, TEAM4-◉)"},
			17: []string{"@u08 (S2-◉, TEAM2-◉)", "@u13 (S1-◉, TEAM3-◉)"},
			18: []string{"@u11 (S1-◉, TEAM3-◉)"},
			19: []string{"@u05 (S2-◉, TEAM1-◉)", "@u10 (S1-◉, TEAM2-◉)", "@u15 (S2-◉, TEAM4-◉)"},
			20: []string{"@u12 (S1-◉, TEAM3-◉)"},
			21: []string{"@u01 (S1-◉, TEAM1-◉)", "@u06 (S2-◉, TEAM2-◉)"},
			22: []string{"@u07 (S2-◉, TEAM2-◉)"},
			23: []string{"@u09 (S1-◉, TEAM2-◉)", "@u18 (S2-◉, TEAM4-◉)"},
			25: []string{"@u16 (S2-◉, TEAM4-◉)"},
			26: []string{"@u03 (S1-◉, TEAM1-◉)", "@u17 (S2-◉, TEAM4-◉)"},
			27: []string{"@u02 (S1-◉, TEAM1-◉)"},
		},
	)
}
