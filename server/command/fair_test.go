// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestFairSimple(t *testing.T) {
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

// TODO TestFairSET is effectively failing; for some reason preferring the
// "MOBILE" team. Need to do some careful analisys, the skills structure may
// indeed be skewed towards non-mobile, but would think randomization would
// help? Need a better "fairness" test harness anyway.
//
// TODO Profile it, it is embarassingly slow
//
// Also, this test relies on the default Limit need selection function,
// pickRequireWeightedRandom. Consider parameterizing.
func TestFairSET(t *testing.T) {
	// 50 shifts, seeds [1-10), no extra fuzz
	multiseedBench(t, "SET", 1, 10, 50, 0, setupUseCaseSET,
		nil,
		map[int][]string{
			60: []string{
				"@fs1lead (TEAMFS1-◉, lead-◈, server-◈, webapp-▣)",
			},
			62: []string{
				"@perfsde2-1 (TEAMPERF-◉, server-◈)",
			},
			64: []string{
				"@perfsde1-1 (TEAMPERF-◉, server-▣, webapp-◉)",
				"@perfsde3-1 (TEAMPERF-◉, server-◈◈)",
			},
			65: []string{
				"@perfsde3-2 (TEAMPERF-◉, server-◈, webapp-▣)",
				"@srelead (TEAMSRE-◉, lead-▣, server-◉, sre-◈, webapp-◉)"},
			66: []string{
				"@fs3sde1-3 (TEAMFS3-◉, server-◉, webapp-◉)",
			},
			68: []string{
				"@fs3sde1-1 (TEAMFS3-◉, webapp-▣)",
				"@fs3sde1-2 (TEAMFS3-◉, server-◉, webapp-▣)",
				"@fs3sde2-2 (TEAMFS3-◉, server-◉, webapp-◈)",
			},
			69: []string{
				"@fs3sde2-1 (TEAMFS3-◉, server-◈, webapp-▣)",
			},
			71: []string{
				"@sre1-1 (TEAMSRE-◉, server-◉, sre-◉)",
			},
			73: []string{
				"@sre2-1 (TEAMSRE-◉, sre-◈, webapp-◉)",
			},
			75: []string{
				"@fs1sde1-2 (TEAMFS1-◉, server-▣, webapp-◉)",
			},
			76: []string{
				"@fs1sde1-1 (TEAMFS1-◉, webapp-▣)",
				"@fs2sde1-2 (TEAMFS2-◉, server-◉, webapp-▣)",
			},
			77: []string{
				"@fs1sde2-2 (TEAMFS1-◉, server-▣, webapp-◉)",
				"@fs2sde1-1 (TEAMFS2-◉, webapp-▣)",
			},
			78: []string{
				"@fs2sde2-1 (TEAMFS2-◉, server-◈, webapp-▣)",
			},
			79: []string{
				"@fs2sde1-3 (TEAMFS2-◉, server-▣, webapp-◉)",
			},
			80: []string{
				"@perflead (TEAMPERF-◉, lead-▣, server-◈, webapp-▣)",
			},
			81: []string{
				"@mobilelead (TEAMMOBILE-◉, lead-◈, mobile-◈, webapp-◈)",
			},
			82: []string{
				"@fs2lead (TEAMFS2-◉, lead-◈, server-◈, webapp-◈)",
				"@fs3lead (TEAMFS3-◉, lead-◈, mobile-▣, server-◉, webapp-◈)",
			},
			83: []string{
				"@fs1sde2-1 (TEAMFS1-◉, mobile-◉, server-▣, webapp-▣)",
				"@perfsde3-3 (TEAMPERF-◉, mobile-◈, webapp-◈)",
			},
			86: []string{
				"@mobilesde2-1 (TEAMMOBILE-◉, mobile-◈, webapp-▣)",
				"@sre2-2 (TEAMSRE-◉, mobile-◉, sre-◈)",
			},
			91: []string{
				"@mobilesde2-2 (TEAMMOBILE-◉, mobile-▣, webapp-▣)",
			},
			95: []string{
				"@mobilesde1-1 (TEAMMOBILE-◉, mobile-▣)",
			},
		},
	)
}

func multiseedBench(t testing.TB, rotationID string, lowSeed, highSeed int64, numShifts, fuzz int, setup string, expectedRetries map[int]int, expectedOut map[int][]string) {
	retries := map[int]int{}
	served := map[string]int{}
	for seed := lowSeed; seed < highSeed; seed++ {
		ctrl, SL := defaultEnv(t)
		defer ctrl.Finish()
		mustRunMulti(t, SL, setup)
		mustRun(t, SL, fmt.Sprintf(`/lotto rotation set fill %s --fuzz 0 --seed %v`, rotationID, seed))
		for n := 0; n < numShifts; n++ {
			mustRun(t, SL, fmt.Sprintf(`/lotto task new shift %s --number %v`, rotationID, n))
			retry := 0
			for {
				taskID := types.ID(fmt.Sprintf(`%s#%v`, rotationID, n))
				_, err := run(t, SL, fmt.Sprintf(`/lotto task fill %v`, taskID))
				if err == nil {
					mustRun(t, SL, fmt.Sprintf(`/lotto task schedule %v`, taskID))
					task, err1 := SL.LoadTask(taskID)
					require.NoError(t, err1)
					for _, user := range task.Users.AsArray() {
						key := user.MarkdownWithSkills().String()
						served[key]++
					}
					break
				}

				retry++
				if retry == 3 {
					require.FailNow(t, "3 retries exceeded: "+err.Error())
				}
				retries[n] = retry
				mustRun(t, SL, fmt.Sprintf(`/lotto rotation set fill %s --seed=%v`, rotationID, time.Now().Unix()))
			}
		}
	}

	if expectedRetries == nil {
		expectedRetries = map[int]int{}
	}
	require.Equal(t, expectedRetries, retries)

	out := map[int][]string{}
	for usermd, count := range served {
		out[count] = append(out[count], usermd)
		sort.Strings(out[count])
	}
	require.EqualValues(t, expectedOut, out)
}
