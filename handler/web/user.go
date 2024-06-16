package web

import (
	"context"
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/services"
)

type UserHandler struct {
	*services.UserService
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return config.Render(w, "login", map[string]any{})
	case http.MethodPost:

		//email := r.FormValue("email")
		//password := r.FormValue("password")

		//statusCode, response := h.LoginService(w, r)
		//return config.WriteJSON(w, statusCode, response)

		user := &models.User{}
		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
		r = r.WithContext(ctx)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	return nil
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	data := map[string]any{}
	return config.Render(w, "register", data)
}

func (h *UserHandler) HomeHandler(w http.ResponseWriter, _ *http.Request) error {
	return config.Render(w, "home", map[string]any{})
}

func (h *UserHandler) ListHandler(w http.ResponseWriter, _ *http.Request) error {
	return config.Render(w, "schedule", map[string]any{})
}

func (h *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "profile", map[string]any{})
}
