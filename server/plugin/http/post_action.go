// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package http

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-msoffice/server/api"
)

func (h *Handler) preprocessAction(w http.ResponseWriter, req *http.Request) (api.API, string) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	if request == nil {
		// h.LogWarn("failed to decode PostActionIntegrationRequest")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return nil, ""
	}
	option, _ := request.Context["selected_option"].(string)

	return api.FromContext(req.Context()), option
}

func (h *Handler) actionRespond(w http.ResponseWriter, req *http.Request) {
	_, _ = h.preprocessAction(w, req)
}
