package web

import (
	"net/http"

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
				setAuthCookie(w, token)
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

func (h *UserHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.PassUpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/profile")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	clearAuthCookie(w)
	// htmx submits this as a POST, so the redirect has to be handed to it as a header
	w.Header().Set("HX-Redirect", "/login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}
