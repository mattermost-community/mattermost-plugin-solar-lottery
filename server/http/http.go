// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
)

// Handler is an http.Handler for all plugin HTTP endpoints
type Handler struct {
	// Logger utils.Logger
	*mux.Router
}

// InitRouter initializes the router.
func NewHandler() *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	apiRouter := h.Router.PathPrefix(config.PathAPI).Subrouter()
	apiRouter.HandleFunc("/authorized", h.apiGetAuthorized).Methods("GET")

	h.Router.Handle("{anything:.*}", http.NotFoundHandler())
	return h
}

func (h *Handler) jsonError(w http.ResponseWriter, statusCode int, summary string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(struct {
		Error   string `json:"error"`
		Summary string `json:"details"`
	}{
		Summary: summary,
		// Summary:   "An internal error has occurred. Check app server logs for details.",
		Error: err.Error(),
	})
	_, _ = w.Write(b)
}

func (h *Handler) internalServerError(w http.ResponseWriter, err error) {
	h.jsonError(w, http.StatusInternalServerError, "An internal error has occurred. Check app server logs for details.", err)
}

func (h *Handler) badRequest(w http.ResponseWriter, err error) {
	h.jsonError(w, http.StatusBadRequest, "Invalid request.", err)
}

func (h *Handler) notFound(w http.ResponseWriter, err error) {
	h.jsonError(w, http.StatusNotFound, "Not found.", err)
}
