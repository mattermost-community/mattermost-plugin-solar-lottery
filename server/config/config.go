package config

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"

// StoredConfig represents the data stored in and managed with the Mattermost
// config.
type StoredConfig struct {
	bot.BotConfig
}

func (c StoredConfig) ToStorableConfig(configMap map[string]interface{}) map[string]interface{} {
	return c.BotConfig.ToStorableConfig(nil)
}

// Config represents the the metadata handed to all request runners (command,
// http).
type Config struct {
	StoredConfig

	BotUserID              string
	BuildDate              string
	BuildHash              string
	BuildHashShort         string
	MattermostSiteHostname string
	MattermostSiteURL      string
	PluginID               string
	PluginURL              string
	PluginURLPath          string
	PluginVersion          string
}
