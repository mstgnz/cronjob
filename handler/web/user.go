package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/load"
	"github.com/mstgnz/cronjob/services"
)

type UserHandler struct {
	*services.UserService
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return load.Render(w, r, "login", map[string]any{})
	case http.MethodPost:
		code, response := h.LoginService(w, r)
		if response.Status && code == http.StatusOK {
			user, ok := response.Data["user"].(*models.User)
			token, ok1 := response.Data["token"].(string)
			if ok && ok1 && user.ID > 0 {
				w.Header().Set("HX-Redirect", "/")
				http.SetCookie(w, &http.Cookie{
					Name:    "Authorization",
					Value:   strings.Join([]string{"Bearer", token}, " "),
					Expires: time.Now().Add(12 * time.Hour),
				})
			}
		}
		_, _ = w.Write([]byte(response.Message))
	}
	return nil
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return load.Render(w, r, "register", map[string]any{})
	case http.MethodPost:
		code, response := h.RegisterService(w, r)
		if response.Status && code == http.StatusCreated {
			user, ok := response.Data["user"].(*models.User)
			token, ok1 := response.Data["token"].(string)
			if ok && ok1 && user.ID > 0 {
				w.Header().Set("HX-Redirect", "/")
				http.SetCookie(w, &http.Cookie{
					Name:    "Authorization",
					Value:   strings.Join([]string{"Bearer", token}, " "),
					Expires: time.Now().Add(12 * time.Hour),
				})
			}
		}
		_, _ = w.Write([]byte(response.Message))
	}
	return nil
}

func (h *UserHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return load.Render(w, r, "home", map[string]any{})
}

func (h *UserHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	return load.Render(w, r, "schedule", map[string]any{})
}

func (h *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	return load.Render(w, r, "profile", map[string]any{})
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
