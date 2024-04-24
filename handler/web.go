package handler

import (
	"net/http"
	"time"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/service"
)

type Web struct {
	service.Api
}

func (wb *Web) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return config.Render(w, "login", map[string]any{})
	case http.MethodPost:
		return config.WriteJSON(w, http.StatusAccepted, config.Response{Status: true, Message: "login", Data: r.Body})
	}
	return nil
}

func (wb *Web) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	data := map[string]any{}
	if r.Method == http.MethodPost {
		_, response := wb.Api.RegisterService(w, r)
		if response.Status {
			token, _ := response.Data.(string)
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour),
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		}
		data["status"] = response.Status
		data["message"] = response.Message
	}
	return config.Render(w, "register", data)
}

func (wb *Web) HomeHandler(w http.ResponseWriter, _ *http.Request) error {
	return config.Render(w, "home", map[string]any{})
}

func (wb *Web) ListHandler(w http.ResponseWriter, _ *http.Request) error {
	return config.Render(w, "schedule", map[string]any{})
}

func (wb *Web) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "profile", map[string]any{})
}
