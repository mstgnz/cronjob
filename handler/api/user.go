package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type UserHandler struct {
	*services.UserService
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.LoginService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.RegisterService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.ProfileService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.UpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *UserHandler) PassUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.PassUpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.DeleteService(w, r)
	return response.WriteJSON(w, statusCode, result)
}
