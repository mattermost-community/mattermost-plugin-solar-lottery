// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/mock_sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/mattermost/mattermost-server/v5/model"
)

var testConfig = config.Config{
	StoredConfig: &config.StoredConfig{
		BotConfig: bot.BotConfig{},
	},
	BuildConfig: &config.BuildConfig{
		PluginID:       "test-plugin-id",
		PluginVersion:  "test-plugin-version",
		BuildDate:      "test-build-date",
		BuildHash:      "test-build-hash",
		BuildHashShort: "test-build-hash-short",
	},
	BotUserID:              "test-bot-user-id",
	MattermostSiteHostname: "siteurl",
	MattermostSiteURL:      "https://siteurl",
	PluginURL:              "https://pluginurl",
	PluginURLPath:          "test-plugin-path",
}

func testUserTimezone() model.StringMap {
	return model.StringMap{
		"useAutomaticTimezone": "false",
		"manualTimezone":       "America/Los_Angeles",
	}
}

func defaultEnv(t testing.TB) (*gomock.Controller, sl.SL) {
	ctrl := gomock.NewController(t)
	sl, _ := getTestSL(t, ctrl)
	return ctrl, sl
}

func getTestSL(t testing.TB, ctrl *gomock.Controller) (sl.SL, kvstore.Store) {
	return getTestSLWithPoster(t, ctrl, nil)
}

func getTestSLWithPoster(t testing.TB, ctrl *gomock.Controller, poster bot.Poster) (sl.SL, kvstore.Store) {
	pluginAPI := mock_sl.NewMockPluginAPI(ctrl)

	pluginAPI.EXPECT().GetMattermostUser(gomock.Any()).AnyTimes().DoAndReturn(func(id string) (*model.User, error) {
		user := &model.User{
			Id:       id,
			Username: id,
			Timezone: testUserTimezone(),
		}
		return user, nil
	})

	pluginAPI.EXPECT().GetMattermostUserByUsername(gomock.Any()).AnyTimes().DoAndReturn(func(username string) (*model.User, error) {
		user := &model.User{
			Id:       username,
			Username: username,
			Timezone: testUserTimezone(),
		}
		return user, nil
	})

	if poster == nil {
		poster = &bot.NilPoster{}
	}
	serviceSL := &sl.Service{
		PluginAPI: pluginAPI,
		Config:    config.NewTestService(&testConfig),
		TaskFillers: map[types.ID]sl.TaskFiller{
			// queue.Type:        queue.New(bot),
			solarlottery.Type: solarlottery.New(),
			"":                solarlottery.New(), // default
		},
		// Logger: &bot.TestLogger{TB: t},
		Logger: &bot.NilLogger{},
		Poster: poster,
		Store:  kvstore.NewStore(kvstore.NewCacheKVStore(nil)),
	}

	return serviceSL.ActingAs("test-user"), serviceSL.Store
}

func run(t testing.TB, sl sl.SL, cmd string) (md.MD, error) {
	if cmd == "" || cmd[0] == '#' {
		return "", nil
	}
	split := strings.Fields(strings.TrimSpace(cmd))
	require.Greater(t, len(split), 1)
	c := &Command{
		SL:            sl,
		actualTrigger: split[0],
	}

	return c.main(split[1:])
}

func mustRun(t testing.TB, sl sl.SL, cmd string) md.MD {
	out, err := run(t, sl, cmd)
	require.NoError(t, err)
	return out
}

func mustRunMulti(t testing.TB, sl sl.SL, in string) {
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		mustRun(t, sl, strings.TrimSpace(line))
	}
}

func runJSON(t testing.TB, sl sl.SL, cmd string, ref interface{}) (md.MD, error) {
	cmd += " --json"
	outmd, err := run(t, sl, cmd)
	if err != nil {
		return "", err
	}
	out := strings.Trim(strings.TrimSpace(outmd.String()), "`")
	out = strings.TrimPrefix(out, "json\n")
	if ref != nil && out != "" {
		errUnmarshal := json.Unmarshal([]byte(out), ref)
		if errUnmarshal != nil {
			return "", errUnmarshal
		}
	}
	return md.MD(out), nil
}

func mustRunJSON(t testing.TB, sl sl.SL, cmd string, ref interface{}) md.MD {
	out, err := runJSON(t, sl, cmd, ref)
	require.NoError(t, err)
	return out
}

func mustRunTaskCreate(t testing.TB, s sl.SL, cmd string) *sl.Task {
	out := &sl.OutCreateTask{}
	mustRunJSON(t, s, cmd, &out)
	return out.Task
}

func mustRunTask(t testing.TB, s sl.SL, cmd string) *sl.Task {
	out := &sl.Task{}
	mustRunJSON(t, s, cmd, &out)
	return out
}

func mustRunRotation(t testing.TB, s sl.SL, cmd string) *sl.Rotation {
	out := &sl.Rotation{}
	mustRunJSON(t, s, cmd, &out)
	return out
}

func mustRunUser(t testing.TB, s sl.SL, cmd string) *sl.User {
	out := sl.NewUser("")
	mustRunJSON(t, s, cmd, &out)
	return out
}

func mustRunTaskAssign(t testing.TB, s sl.SL, cmd string) *sl.Task {
	out := &sl.OutAssignTask{
		Changed: sl.NewUsers(),
	}
	mustRunJSON(t, s, cmd, &out)
	return out.Task
}

func mustRunUsersQualify(t testing.TB, s sl.SL, cmd string) *sl.Users {
	out := sl.OutQualify{
		Users: sl.NewUsers(),
	}
	mustRunJSON(t, s, cmd, &out)
	return out.Users
}

func mustRunUsersJoin(t testing.TB, s sl.SL, cmd string) *sl.Users {
	out := sl.OutJoinRotation{
		Modified: sl.NewUsers(),
	}
	mustRunJSON(t, s, cmd, &out)
	return out.Modified
}

func mustRunUsers(t testing.TB, s sl.SL, cmd string) *sl.Users {
	out := sl.NewUsers()
	mustRunJSON(t, s, cmd, &out)
	return out
}

func mustRunUsersCalendar(t testing.TB, s sl.SL, cmd string) *sl.Users {
	out := sl.OutCalendar{
		Users: sl.NewUsers(),
	}
	mustRunJSON(t, s, cmd, &out)
	return out.Users
}
