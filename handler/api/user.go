package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type UserHandler struct {
	*services.UserService
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.LoginService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.RegisterService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ProfileService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *UserHandler) PassUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.PassUpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
