// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
)

const (
	PathAPI        = "/api/v1"
	PathPostAction = "/action"
	PathRespond    = "/respond"
)

// Handler is an http.Handler for all plugin HTTP endpoints
type Service struct {
	config config.Service
	*mux.Router
	sl sl.Service
}

func NewService(config config.Service, router *mux.Router, sl sl.Service) *Service {
	s := &Service{
		Router: router,
		config: config,
		sl:     sl,
	}
	apiRouter := s.Router.PathPrefix(PathAPI).Subrouter()
	apiRouter.HandleFunc("/authorized", s.apiGetAuthorized).Methods("GET")
	apiRouter.HandleFunc("/execute_command", s.executeCommand).Methods("POST")

	return s
}

func (s *Service) handleErrorWithCode(w http.ResponseWriter, code int, errTitle string, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(struct {
		Error   string `json:"error"`
		Details string `json:"details"`
	}{
		Error:   errTitle,
		Details: err.Error(),
	})
	_, _ = w.Write(b)
}
