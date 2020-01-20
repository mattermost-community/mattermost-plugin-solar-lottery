// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package plugin

import (
	"math/rand"
	gohttp "net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/command"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/http"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
)

type Plugin struct {
	plugin.MattermostPlugin
	configLock *sync.RWMutex
	config     *config.Config

	httpHandler *http.Handler
	// notificationHandler api.NotificationHandler

	Templates map[string]*template.Template
}

func NewWithConfig(conf *config.Config) *Plugin {
	return &Plugin{
		configLock: &sync.RWMutex{},
		config:     conf,
	}
}

func (p *Plugin) OnActivate() error {
	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    config.BotUserName,
		DisplayName: config.BotDisplayName,
		Description: config.BotDescription,
	}, plugin.ProfileImagePath("assets/profile.png"))
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot account")
	}
	p.config.BotUserID = botUserID

	// Templates
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "couldn't get bundle path")
	}
	err = p.loadTemplates(bundlePath)
	if err != nil {
		return err
	}

	p.httpHandler = http.NewHandler()
	// p.notificationHandler = api.NewNotificationHandler(p.newAPIConfig())

	command.Register(p.API.RegisterCommand)

	rand.Seed(time.Now().UnixNano())

	p.API.LogInfo(p.config.PluginID + " activated")
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	conf := p.getConfig()
	stored := config.StoredConfig{}
	err := p.API.LoadPluginConfiguration(&stored)
	if err != nil {
		return errors.WithMessage(err, "failed to load plugin configuration")
	}

	mattermostSiteURL := p.API.GetConfig().ServiceSettings.SiteURL
	if mattermostSiteURL == nil {
		return errors.New("plugin requires Mattermost Site URL to be set")
	}
	mattermostURL, err := url.Parse(*mattermostSiteURL)
	if err != nil {
		return err
	}
	pluginURLPath := "/plugins/" + conf.PluginID
	pluginURL := strings.TrimRight(*mattermostSiteURL, "/") + pluginURLPath

	p.updateConfig(func(c *config.Config) {
		c.StoredConfig = stored
		c.MattermostSiteURL = *mattermostSiteURL
		c.MattermostSiteHostname = mattermostURL.Hostname()
		c.PluginURL = pluginURL
		c.PluginURLPath = pluginURLPath
	})

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	wasDemo := p.executeDemoCommand(c, args)
	if wasDemo {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Demo done",
		}, nil
	}

	apiconf := p.newAPIConfig()
	command := command.Command{
		Context:   c,
		Args:      args,
		ChannelID: args.ChannelId,
		Config:    apiconf.Config,
		API:       api.New(apiconf, args.UserId),
	}

	out, _ := command.Handle()
	p.SendEphemeralPost(args.ChannelId, args.UserId, out)
	return &model.CommandResponse{}, nil
}

func (p *Plugin) ServeHTTP(pc *plugin.Context, w gohttp.ResponseWriter, req *gohttp.Request) {
	apiconf := p.newAPIConfig()
	mattermostUserID := req.Header.Get("Mattermost-User-ID")
	ctx := req.Context()
	ctx = api.Context(ctx, api.New(apiconf, mattermostUserID))
	ctx = config.Context(ctx, apiconf.Config)

	p.httpHandler.ServeHTTP(w, req.WithContext(ctx))
}

func (p *Plugin) getConfig() *config.Config {
	p.configLock.RLock()
	defer p.configLock.RUnlock()
	return &(*p.config)
}

func (p *Plugin) updateConfig(f func(*config.Config)) config.Config {
	p.configLock.Lock()
	defer p.configLock.Unlock()

	f(p.config)
	return *p.config
}

func (p *Plugin) newAPIConfig() api.Config {
	conf := p.getConfig()
	bot := bot.NewBot(p.API, conf.BotUserID).WithConfig(conf.BotConfig)
	store := store.NewPluginStore(p.API, bot)

	return api.Config{
		Config: conf,
		Dependencies: &api.Dependencies{
			Autofillers: map[string]api.Autofiller{
				"":                solarlottery.New(bot), // default
				solarlottery.Type: solarlottery.New(bot),
			},
			RotationStore: store,
			SkillsStore:   store,
			UserStore:     store,
			ShiftStore:    store,
			Logger:        bot,
			Poster:        bot,
			PluginAPI:     p,
		},
	}
}

func (p *Plugin) loadTemplates(bundlePath string) error {
	if p.Templates != nil {
		return nil
	}

	templatesPath := filepath.Join(bundlePath, "assets", "templates")
	templates := make(map[string]*template.Template)
	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		template, err := template.ParseFiles(path)
		if err != nil {
			return nil
		}
		key := path[len(templatesPath):]
		templates[key] = template
		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "OnActivate/loadTemplates failed")
	}
	p.Templates = templates
	return nil
}
