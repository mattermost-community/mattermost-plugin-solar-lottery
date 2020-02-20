// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
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
}

func NewService(config config.Service, router *mux.Router) *Service {
	s := &Service{
		Router: router,
		config: config,
	}
	apiRouter := s.Router.PathPrefix(PathAPI).Subrouter()
	apiRouter.HandleFunc("/authorized", s.apiGetAuthorized).Methods("GET")
	return s
}
