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
