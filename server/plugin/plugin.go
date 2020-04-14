// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package plugin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/command"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/constants"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl/filler/solarlottery"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Plugin struct {
	plugin.MattermostPlugin

	bot    bot.Bot
	sl     sl.Service
	api    *api.Service
	config config.Service

	botUserID string
}

func New(build *config.BuildConfig) *Plugin {
	p := &Plugin{}
	p.config = config.NewService(build, p)
	return p
}

func (p *Plugin) OnActivate() error {
	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    constants.BotUserName,
		DisplayName: constants.BotDisplayName,
		Description: constants.BotDescription,
	}, plugin.ProfileImagePath("assets/profile.png"))
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot account")
	}
	p.botUserID = botUserID
	p.bot = bot.NewBot(p.API, botUserID)

	p.sl = sl.Service{
		PluginAPI: p,
		Config:    p.config,
		TaskFillers: map[types.ID]sl.TaskFiller{
			// queue.Type:        queue.New(bot),
			solarlottery.Type: solarlottery.New(),
			"":                solarlottery.New(), // default
		},
		Logger: p.bot,
		Poster: p.bot,
		Store:  kvstore.NewStore(kvstore.NewPluginStore(p.API)),
	}

	router := &mux.Router{}
	p.api = api.NewService(p.config, router)
	router.Handle("{anything:.*}", http.NotFoundHandler())

	command.Register(p.API.RegisterCommand)
	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	return p.config.Refresh()
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	wasDemo := p.executeDemoCommand(c, args)
	if wasDemo {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Demo done",
		}, nil
	}

	command := command.Command{
		Context:   c,
		Args:      args,
		ChannelID: args.ChannelId,
		SL:        p.sl.ActingAs(types.ID(args.UserId)),
	}

	out, _ := command.Handle()
	p.SendEphemeralPost(args.ChannelId, args.UserId, out.String())
	return &model.CommandResponse{}, nil
}

func (p *Plugin) ServeHTTP(pc *plugin.Context, w http.ResponseWriter, req *http.Request) {
	p.api.ServeHTTP(w, req)
}
