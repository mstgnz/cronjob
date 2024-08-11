package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type NotifyEmailHandler struct {
	*services.NotifyEmailService
}

func (h *NotifyEmailHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.ListService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *NotifyEmailHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.CreateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *NotifyEmailHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.UpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *NotifyEmailHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.DeleteService(w, r)
	return response.WriteJSON(w, statusCode, result)
}
