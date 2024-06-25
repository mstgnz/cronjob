package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type GroupHandler struct {
	*services.GroupService
}

func (h *GroupHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ListService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *GroupHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.CreateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *GroupHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *GroupHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
