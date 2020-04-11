// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func TestTaskUnassign(t *testing.T) {
	for _, tc := range []struct {
		name          string
		assigned      string
		transitions   []string
		unassign      string
		force         bool
		expectError   bool
		expectChanged []string
	}{
		{
			name:          "happy",
			assigned:      "@test-user1 @test-user2",
			unassign:      "@test-user1 @test-user2",
			expectChanged: []string{"test-user1", "test-user2"},
		}, {
			name:          "happy force scheduled",
			assigned:      "@test-user1 @test-user2",
			transitions:   []string{"schedule"},
			unassign:      "@test-user1 ",
			force:         true,
			expectChanged: []string{"test-user1"},
		}, {
			name:          "happy force started",
			assigned:      "@test-user1 @test-user2",
			transitions:   []string{"schedule", "start"},
			unassign:      "@test-user1 ",
			force:         true,
			expectChanged: []string{"test-user1"},
		}, {
			name:        "fail scheduled",
			assigned:    "@test-user1 @test-user2",
			transitions: []string{"schedule"},
			unassign:    "@test-user1 ",
			expectError: true,
		}, {
			name:        "fail started",
			assigned:    "@test-user1 @test-user2",
			transitions: []string{"schedule", "start"},
			unassign:    "@test-user1 ",
			expectError: true,
		}, {
			name:        "fail finished",
			assigned:    "@test-user1 @test-user2",
			transitions: []string{"schedule", "start", "finish"},
			unassign:    "@test-user1 ",
			expectError: true,
		}, {
			name:        "fail force finished",
			assigned:    "@test-user1 @test-user2",
			transitions: []string{"schedule", "start", "finish"},
			unassign:    "@test-user1 ",
			force:       true,
			expectError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			SL, _ := getTestSL(t, ctrl)

			err := runCommands(t, SL, `
				/lotto rotation new test-rotation
				/lotto rotation param ticket test-rotation
				/lotto task new ticket test-rotation --summary test-summary1
			`)
			require.NoError(t, err)
			err = runCommands(t, SL, "/lotto task assign test-rotation#1 "+tc.assigned)
			require.NoError(t, err)

			for _, transition := range tc.transitions {
				err = runCommands(t, SL, "/lotto task "+transition+" test-rotation#1")
				require.NoError(t, err)
			}

			out := &sl.OutAssignTask{
				Changed: sl.NewUsers(),
			}
			force := " "
			if tc.force {
				force = force + "--force"
			}
			_, err = runJSONCommand(t, SL, "/lotto task unassign test-rotation#1 "+tc.unassign+force, &out)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, types.ID("test-rotation#1"), out.Task.TaskID)
			require.Equal(t, tc.expectChanged, out.Changed.TestIDs())
		})
	}
}
