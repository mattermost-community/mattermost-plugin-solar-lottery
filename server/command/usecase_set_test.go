// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

const setupUseCaseSET = `
/lotto rotation new SET --beginning=2019-01-16 --period=monthly
/lotto rotation set task SET --grace=300h 

/lotto rotation set require SET -s any --count 5
/lotto rotation set require SET -s webapp --count 1
/lotto rotation set require SET -s server --count 1
/lotto rotation set require SET -s mobile --count 1
/lotto rotation set require SET -s lead --count 1

/lotto rotation set limit SET -s any --count 6
/lotto rotation set limit SET -s lead --count 1
/lotto rotation set limit SET -s TEAMFS1 --count 1
/lotto rotation set limit SET -s TEAMFS2 --count 1
/lotto rotation set limit SET -s TEAMFS3 --count 1
/lotto rotation set limit SET -s TEAMPERF --count 1
/lotto rotation set limit SET -s TEAMSRE --count 1
/lotto rotation set limit SET -s TEAMMOBILE --count 1

/lotto user qualify @fs1lead       -s TEAMFS1,lead-3,server-3,webapp-2
/lotto user qualify @fs1sde2-1     -s TEAMFS1,server-2,webapp-2,mobile-1
/lotto user qualify @fs1sde2-2     -s TEAMFS1,server-2,webapp-1
/lotto user qualify @fs1sde1-1     -s TEAMFS1,webapp-2
/lotto user qualify @fs1sde1-2     -s TEAMFS1,server-2,webapp-1

/lotto user qualify @fs2lead       -s TEAMFS2,lead-3,server-3,webapp-3
/lotto user qualify @fs2sde2-1     -s TEAMFS2,server-3,webapp-2
/lotto user qualify @fs2sde1-1     -s TEAMFS2,webapp-2
/lotto user qualify @fs2sde1-2     -s TEAMFS2,server-1,webapp-2
/lotto user qualify @fs2sde1-3     -s TEAMFS2,server-2,webapp-1

/lotto user qualify @fs3lead       -s TEAMFS3,lead-3,server-1,webapp-3,mobile-2
/lotto user qualify @fs3sde2-1     -s TEAMFS3,server-3,webapp-2
/lotto user qualify @fs3sde2-2     -s TEAMFS3,server-1,webapp-3
/lotto user qualify @fs3sde1-1     -s TEAMFS3,webapp-2
/lotto user qualify @fs3sde1-2     -s TEAMFS3,server-1,webapp-2
/lotto user qualify @fs3sde1-3     -s TEAMFS3,server-1,webapp-1

/lotto user qualify @mobilelead    -s TEAMMOBILE,lead-3,webapp-3,mobile-3
/lotto user qualify @mobilesde2-1  -s TEAMMOBILE,mobile-3,webapp-2
/lotto user qualify @mobilesde2-2  -s TEAMMOBILE,mobile-2,webapp-2
/lotto user qualify @mobilesde1-1  -s TEAMMOBILE,mobile-2

/lotto user qualify @perflead      -s TEAMPERF,lead-2,server-3,webapp-2
/lotto user qualify @perfsde3-1    -s TEAMPERF,server-4
/lotto user qualify @perfsde3-2    -s TEAMPERF,server-3,webapp-2
/lotto user qualify @perfsde3-3    -s TEAMPERF,webapp-3,mobile-3
/lotto user qualify @perfsde2-1    -s TEAMPERF,server-3
/lotto user qualify @perfsde1-1    -s TEAMPERF,server-2,webapp-1  

/lotto user qualify @srelead       -s TEAMSRE,lead-2,server-1,webapp-1,sre-3
/lotto user qualify @sre2-1        -s TEAMSRE,sre-3,webapp-1
/lotto user qualify @sre2-2        -s TEAMSRE,sre-3,mobile-1
/lotto user qualify @sre1-1        -s TEAMSRE,sre-1,server-1
/lotto user qualify @sre1-1        -s TEAMSRE,sre-1

/lotto user join SET @fs1lead @fs1sde2-1 @fs1sde2-2 @fs1sde1-1 @fs1sde1-2 
/lotto user join SET @fs2lead @fs2sde2-1 @fs2sde1-1 @fs2sde1-2 @fs2sde1-3  
/lotto user join SET @fs3lead @fs3sde2-1 @fs3sde2-2 @fs3sde1-1 @fs3sde1-2 @fs3sde1-3  
/lotto user join SET @mobilelead @mobilesde2-1 @mobilesde2-2 @mobilesde1-1
/lotto user join SET @perflead @perfsde3-1 @perfsde3-2 @perfsde3-3 @perfsde2-1 @perfsde1-1 
/lotto user join SET @srelead @sre2-1 @sre2-2 @sre1-1 @sre1-1     
`

func TestUseCaseSETHappy(t *testing.T) {
	ctrl, SL := defaultEnv(t)
	defer ctrl.Finish()
	mustRunMulti(t, SL, setupUseCaseSET)

	mustRun(t, SL,
		`/lotto rotation set fill SET --seed `+strconv.FormatInt(intNoValue, 10))

	checkShift := func(n int, start string, durHours int, assigned []string) {
		task := mustRunTaskCreate(t, SL,
			fmt.Sprintf(`/lotto task new shift SET --number %v`, n))
		require.Equal(t, "UTC", task.ExpectedStart.Location().String())
		require.Equal(t, types.ID(fmt.Sprintf("SET#%v", n)), task.TaskID)
		require.Equal(t, start+"T08:00", task.ExpectedStart.String())
		require.Equal(t, (time.Duration(durHours) * time.Hour).String(), task.ExpectedDuration.String())
		require.Equal(t, []types.ID{"any", "lead-◉", "mobile-◉", "server-◉", "webapp-◉"}, task.Require.IDs())
		require.Equal(t, map[types.ID]int64{"any": 5, "lead-◉": 1, "mobile-◉": 1, "server-◉": 1, "webapp-◉": 1}, task.Require.TestAsMap())
		require.Equal(t, map[types.ID]int64{"TEAMFS1-◉": 1, "TEAMFS2-◉": 1, "TEAMFS3-◉": 1, "TEAMMOBILE-◉": 1, "TEAMPERF-◉": 1, "TEAMSRE-◉": 1, "any": 6, "lead-◉": 1}, task.Limit.TestAsMap())

		out := mustRun(t, SL,
			fmt.Sprintf(`/lotto task fill SET#%v`, n))
		require.Equal(t, md.Markdownf("Auto-assigned %s to ticket SET#%v", strings.Join(assigned, ", "), n), out)
	}

	checkShift(0, "2019-01-16", 744, []string{
		"@mobilelead (TEAMMOBILE-◉, lead-◈, mobile-◈, webapp-◈)",
		"@perfsde2-1 (TEAMPERF-◉, server-◈)",
		"@fs3sde1-3 (TEAMFS3-◉, server-◉, webapp-◉)",
		"@sre1-1 (TEAMSRE-◉, server-◉, sre-◉)",
		"@fs2sde2-1 (TEAMFS2-◉, server-◈, webapp-▣)",
	})

	// Make sure the shift us added to the users' calendars (check karen.austin)
	u := mustRunUser(t, SL,
		`/lotto user show @mobilelead`)
	require.Equal(t, types.ID("mobilelead"), u.MattermostUserID)
	require.Equal(t, []*sl.Unavailable{
		{Interval: types.MustParseInterval("2019-01-16T08:00", "2019-02-16T08:00"), Reason: "task", TaskID: "SET#0", RotationID: "SET"},
		{Interval: types.MustParseInterval("2019-02-16T08:00", "2019-02-28T20:00"), Reason: "grace", TaskID: "SET#0", RotationID: "SET"},
	}, u.Calendar)

	// Shift #1 should've included @fs3lead, but let's make @fs3lead unavailable
	mustRun(t, SL,
		`/lotto user unavailable @fs3lead --start 2019-02-20T11:00 --finish 2019-03-30T09:30`)
	u = mustRunUser(t, SL,
		`/lotto user show @fs3lead`)
	require.Equal(t, []*sl.Unavailable{{Interval: types.MustParseInterval("2019-02-20T19:00", "2019-03-30T16:30"), Reason: "personal"}}, u.Calendar)

	// Note @fs3lead is not selected
	checkShift(1, "2019-02-16", 672, []string{
		"@perflead (TEAMPERF-◉, lead-▣, server-◈, webapp-▣)",
		"@fs1sde2-1 (TEAMFS1-◉, mobile-◉, server-▣, webapp-▣)",
		"@fs3sde1-2 (TEAMFS3-◉, server-◉, webapp-▣)",
		"@mobilesde2-2 (TEAMMOBILE-◉, mobile-▣, webapp-▣)",
		"@fs2sde1-1 (TEAMFS2-◉, webapp-▣)",
	})
	checkShift(2, "2019-03-16", 744, []string{
		"@mobilelead (TEAMMOBILE-◉, lead-◈, mobile-◈, webapp-◈)",
		"@perfsde2-1 (TEAMPERF-◉, server-◈)",
		"@fs1sde2-2 (TEAMFS1-◉, server-▣, webapp-◉)",
		"@fs3sde2-2 (TEAMFS3-◉, server-◉, webapp-◈)",
		"@sre2-2 (TEAMSRE-◉, mobile-◉, sre-◈)",
	})
	// Note @fs3lead is back in rotation, gets picked up
	checkShift(3, "2019-04-16", 720, []string{
		"@fs3lead (TEAMFS3-◉, lead-◈, mobile-▣, server-◉, webapp-◈)",
		"@perfsde3-3 (TEAMPERF-◉, mobile-◈, webapp-◈)",
		"@fs1sde1-1 (TEAMFS1-◉, webapp-▣)",
		"@mobilesde1-1 (TEAMMOBILE-◉, mobile-▣)",
		"@fs2sde2-1 (TEAMFS2-◉, server-◈, webapp-▣)",
	})
	checkShift(4, "2019-05-16", 744, []string{
		"@fs2lead (TEAMFS2-◉, lead-◈, server-◈, webapp-◈)",
		"@mobilesde2-2 (TEAMMOBILE-◉, mobile-▣, webapp-▣)",
		"@perfsde1-1 (TEAMPERF-◉, server-▣, webapp-◉)",
		"@fs1sde2-1 (TEAMFS1-◉, mobile-◉, server-▣, webapp-▣)",
		"@fs3sde1-1 (TEAMFS3-◉, webapp-▣)",
	})
	checkShift(5, "2019-06-16", 720, []string{
		"@perflead (TEAMPERF-◉, lead-▣, server-◈, webapp-▣)",
		"@mobilesde2-1 (TEAMMOBILE-◉, mobile-◈, webapp-▣)",
		"@fs3sde1-2 (TEAMFS3-◉, server-◉, webapp-▣)",
		"@fs2sde1-3 (TEAMFS2-◉, server-▣, webapp-◉)",
		"@sre1-1 (TEAMSRE-◉, server-◉, sre-◉)",
	})
	checkShift(6, "2019-07-16", 744, []string{
		"@srelead (TEAMSRE-◉, lead-▣, server-◉, sre-◈, webapp-◉)",
		"@fs1sde2-1 (TEAMFS1-◉, mobile-◉, server-▣, webapp-▣)",
		"@fs3sde1-3 (TEAMFS3-◉, server-◉, webapp-◉)",
		"@fs2sde1-2 (TEAMFS2-◉, server-◉, webapp-▣)",
		"@perfsde3-1 (TEAMPERF-◉, server-◈◈)",
	})
}

// TODO TestUseCaseSETFair is effectively failing; for some reason preferring
// the "MOBILE" team. Need to do some careful analisys, the skills structure may
// indeed be skewed towards non-mobile, but would think randomization would
// help? Need a better "fairness" test harness anyway.
//
// TODO Profile it, it is embarassingly slow
//
// Also, this test relies on the default Limit need selection function,
// pickRequireWeightedRandom. Consider parameterizing.
func TestUseCaseSETFair(t *testing.T) {
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
