package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
		code, response := h.LoginService(w, r)
		if response.Status && code == http.StatusOK {
			user, ok := response.Data["user"].(*models.User)
			if ok && user.ID > 0 {
				http.SetCookie(w, &http.Cookie{
					Name:    "auth",
					Value:   strconv.Itoa(user.ID),
					Expires: time.Now().Add(12 * time.Hour),
				})
			}
		}
		json.NewEncoder(w).Encode(response)
		return nil
	default:
		json.NewEncoder(w).Encode(map[string]any{"status": false, "message": "not supported request", "data": nil})
		return nil
	}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return config.Render(w, "register", map[string]any{})
	case http.MethodPost:
		code, response := h.RegisterService(w, r)
		if response.Status && code == http.StatusCreated {
			user, ok := response.Data["user"].(*models.User)
			if ok && user.ID > 0 {
				http.SetCookie(w, &http.Cookie{
					Name:    "auth",
					Value:   strconv.Itoa(user.ID),
					Expires: time.Now().Add(12 * time.Hour),
				})
			}
		}
		json.NewEncoder(w).Encode(response)
		return nil
	default:
		json.NewEncoder(w).Encode(map[string]any{"status": false, "message": "not supported request", "data": nil})
		return nil
	}
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

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
