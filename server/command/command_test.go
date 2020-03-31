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

func getTestSL(t testing.TB, ctrl *gomock.Controller) (sl.SL, kvstore.Store) {
	return getTestSLWithPoster(t, ctrl, nil)
}

func getTestSLWithPoster(t testing.TB, ctrl *gomock.Controller, poster bot.Poster) (sl.SL, kvstore.Store) {
	pluginAPI := mock_sl.NewMockPluginAPI(ctrl)

	pluginAPI.EXPECT().GetMattermostUser(gomock.Any()).AnyTimes().DoAndReturn(func(id string) (*model.User, error) {
		user := &model.User{
			Id:       id,
			Username: id + "-username",
			Timezone: testUserTimezone(),
		}
		return user, nil
	})

	pluginAPI.EXPECT().GetMattermostUserByUsername(gomock.Any()).AnyTimes().DoAndReturn(func(username string) (*model.User, error) {
		id := strings.TrimSuffix(username, "-username")
		user := &model.User{
			Id:       id,
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

func runCommand(t testing.TB, sl sl.SL, cmd string) (md.MD, error) {
	if cmd == "" || cmd[0] == '#' {
		return "", nil
	}
	split := strings.Fields(strings.TrimSpace(cmd))
	require.Greater(t, len(split), 1)
	c := &Command{
		SL:            sl,
		actualTrigger: split[0],
	}

	return c.handleCommand(c.commands(), split[1:])
}

func runJSONCommand(t testing.TB, sl sl.SL, cmd string, ref interface{}) (md.MD, error) {
	cmd += " --json"
	outmd, err := runCommand(t, sl, cmd)
	if err != nil {
		return outmd, err
	}
	out := outmd.String()
	out = strings.Trim(strings.TrimSpace(out), "`")
	out = strings.TrimPrefix(out, "json\n")
	if ref != nil && out != "" {
		errUnmarshal := json.Unmarshal([]byte(out), ref)
		if errUnmarshal != nil {
			return "", errUnmarshal
		}
	}
	return md.MD(out), err
}

func runCommands(t testing.TB, sl sl.SL, in string) error {
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		_, err := runCommand(t, sl, strings.TrimSpace(line))
		if err != nil {
			return err
		}
	}
	return nil
}
