package api

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/command"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (s *Service) executeCommand(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	args := model.CommandArgsFromJson(r.Body)

	command := command.Command{
		Context:   &plugin.Context{},
		Args:      args,
		ChannelID: args.ChannelId,
		SL:        s.sl.ActingAs(types.ID(args.UserId)),
	}
	out, err := command.Handle()
	if err == nil {
		s.sl.Logger.Errorf("Error while handling command", err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Error while handling command", err)
	}

	w.Write([]byte(out.String()))
}
