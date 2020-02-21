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
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/mock_sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
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

func getTestSL(t testing.TB, ctrl *gomock.Controller) (sl.SL, kvstore.Store) {
	pluginAPI := mock_sl.NewMockPluginAPI(ctrl)

	pluginAPI.EXPECT().GetMattermostUser(gomock.Eq("test-user")).AnyTimes().Return(
		&model.User{
			Id:       "test-user",
			Username: "test-username",
		}, nil)

	serviceSL := &sl.Service{
		PluginAPI: pluginAPI,
		Config:    config.NewTestService(&testConfig),
		// Autofillers map[string]Autofiller
		Logger: &bot.NilLogger{},
		Poster: &bot.NilPoster{},
		Store:  kvstore.NewStore(kvstore.NewCacheKVStore(nil)),
	}

	return serviceSL.ActingAs("test-user"), serviceSL.Store
}

func runCommand(t testing.TB, sl sl.SL, cmd string) (string, error) {
	if cmd == "" || cmd[0] == '#' {
		return "", nil
	}
	split := strings.Fields(strings.TrimSpace(cmd))
	require.Greater(t, len(split), 1)
	c := &Command{
		SL:         sl,
		subcommand: split[0],
	}

	return c.handleCommand(c.commands(), split[1:])
}

func runJSONCommand(t testing.TB, sl sl.SL, cmd string, ref interface{}) (string, error) {
	cmd += " --json"
	out, err := runCommand(t, sl, cmd)
	out = strings.Trim(strings.TrimSpace(out), "`")
	out = strings.TrimPrefix(out, "json\n")

	if ref != nil && out != "" {
		errUnmarshal := json.Unmarshal([]byte(out), ref)
		if errUnmarshal != nil {
			return "", errUnmarshal
		}
	}
	return out, err
}

func runCommands(t testing.TB, sl sl.SL, in string) {
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		runCommand(t, sl, strings.TrimSpace(line))
	}
}
