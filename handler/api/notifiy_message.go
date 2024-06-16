package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type NotifyMessageHandler struct {
	*services.NotifyMessageService
}

func (h *NotifyMessageHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ListService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *NotifyMessageHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.CreateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *NotifyMessageHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *NotifyMessageHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
