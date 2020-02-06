// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package plugin

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

// This file provides an implementation of sl.PluginAPI
var _ sl.PluginAPI = (*Plugin)(nil)

// IsPluginAdmin returns true if the user is authorized to use the workflow plugin's admin-level APIs/commands.
func (p *Plugin) IsPluginAdmin(mattermostUserID string) (bool, error) {
	user, err := p.API.GetUser(mattermostUserID)
	if err != nil {
		return false, err
	}
	if strings.Contains(user.Roles, "system_admin") {
		return true, nil
	}
	conf := p.getConfig()
	bot := bot.NewBot(p.API, conf.BotUserID).WithConfig(conf.BotConfig)
	return bot.IsUserAdmin(mattermostUserID), nil
}

func (p *Plugin) GetMattermostUserByUsername(mattermostUsername string) (*model.User, error) {
	for strings.HasPrefix(mattermostUsername, "@") {
		mattermostUsername = mattermostUsername[1:]
	}
	u, err := p.API.GetUserByUsername(mattermostUsername)
	if err != nil {
		return nil, err
	}
	if u.DeleteAt != 0 {
		return nil, store.ErrNotFound
	}
	return u, nil
}

func (p *Plugin) GetMattermostUser(mattermostUserID string) (*model.User, error) {
	mmuser, err := p.API.GetUser(mattermostUserID)
	if err != nil {
		return nil, err
	}
	if mmuser.DeleteAt != 0 {
		return nil, store.ErrNotFound
	}
	return mmuser, nil
}

func (p *Plugin) Clean() error {
	appErr := p.API.KVDeleteAll()
	if appErr != nil {
		return appErr
	}
	return nil
}

func (p *Plugin) SendEphemeralPost(channelID, userID, message string) {
	ephemeralPost := &model.Post{
		ChannelId: channelID,
		UserId:    p.getConfig().BotUserID,
		Message:   message,
	}
	_ = p.API.SendEphemeralPost(userID, ephemeralPost)
}

func (p *Plugin) UpdateStoredConfig(f func(*config.Config)) {
	conf := p.updateConfig(f)
	go func() {
		p.API.SavePluginConfig(conf.ToStorableConfig(nil))
	}()
	return
}
