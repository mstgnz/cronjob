package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type Web struct{}

func (h *Web) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return config.Render(w, "login", map[string]any{})
	case http.MethodPost:
		return config.WriteJSON(w, http.StatusAccepted, config.Response{Status: true, Message: "login", Data: r.Body})
	}
	return nil
}

func (h *Web) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	data := map[string]any{}
	return config.Render(w, "register", data)
}

func (h *Web) HomeHandler(w http.ResponseWriter, _ *http.Request) error {
	return config.Render(w, "home", map[string]any{})
}

func (h *Web) ListHandler(w http.ResponseWriter, _ *http.Request) error {
	return config.Render(w, "schedule", map[string]any{})
}

func (h *Web) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "profile", map[string]any{})
}
