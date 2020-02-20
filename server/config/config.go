package config

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

// StoredConfig represents the data stored in and managed with the Mattermost
// config.
type StoredConfig struct {
	bot.BotConfig
}

func (c StoredConfig) Map(onto map[string]interface{}) map[string]interface{} {
	out := c.BotConfig.ToStorable(nil)
	return out
}

type BuildConfig struct {
	PluginID       string
	PluginVersion  string
	BuildDate      string
	BuildHash      string
	BuildHashShort string
}

// Config represents the the metadata handed to all request runners (command,
// http).
type Config struct {
	*StoredConfig
	*BuildConfig

	BotUserID              string
	MattermostSiteHostname string
	MattermostSiteURL      string
	PluginURL              string
	PluginURLPath          string
}
