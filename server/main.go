package main

import (
	mattermost "github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	solarlotteryplugin "github.com/mattermost/mattermost-plugin-solar-lottery/server/plugin"
)

var BuildHash string
var BuildHashShort string
var BuildDate string

func main() {
	mattermost.ClientMain(
		solarlotteryplugin.NewWithConfig(
			&config.Config{
				PluginID:       manifest.ID,
				PluginVersion:  manifest.Version,
				BuildHash:      BuildHash,
				BuildHashShort: BuildHashShort,
				BuildDate:      BuildDate,
			}))
}
